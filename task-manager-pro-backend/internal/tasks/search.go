package tasks

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Context global do Redis
var redisCtx = context.Background()

// Gera hash da query para criar chave única por usuário + busca
func hashQuery(q string) string {
	h := sha1.Sum([]byte(q))
	return hex.EncodeToString(h[:])
}

// ---------------------------------------------------------------------------
// SearchTasks — Busca com Cache + Histórico
// ---------------------------------------------------------------------------
func SearchTasks(userID uint, query string) ([]Task, error) {
	// fallback caso Redis não esteja configurado
	if redisClient == nil || redisClient.Client == nil {
		fmt.Println("Redis não configurado. Usando busca direta no Postgres.")
		filter := TaskFilter{Query: query}
		return ListTasks(userID, filter)
	}

	userIDStr := strconv.Itoa(int(userID))
	queryHash := hashQuery(query)

	cacheKey := "tmpro:search:result:" + userIDStr + ":" + queryHash
	historyKey := "tmpro:search:history:" + userIDStr

	// -----------------------------------------------------------------------
	// 1. Tenta pegar a busca do cache
	// -----------------------------------------------------------------------
	cached, err := redisClient.Client.Get(redisCtx, cacheKey).Result()
	if err == nil && cached != "" {
		var tasks []Task

		if json.Unmarshal([]byte(cached), &tasks) == nil {
			// Atualiza histórico mesmo quando pega do cache
			_ = pushSearchHistory(historyKey, query)
			fmt.Println("Resultado retornado do CACHE!")
			return tasks, nil
		}
	} else if err != nil && err != redis.Nil {
		// erro inesperado do Redis (não derruba o sistema)
		fmt.Println("Erro no Redis GET:", err)
	}

	// -----------------------------------------------------------------------
	// 2. Busca no Postgres usando o ListTasks (já existente)
	// -----------------------------------------------------------------------
	filter := TaskFilter{Query: query}
	tasks, err := ListTasks(userID, filter)
	if err != nil {
		return nil, err
	}

	// -----------------------------------------------------------------------
	// 3. Salva no cache (TTL: 30s)
	// -----------------------------------------------------------------------
	jsonData, _ := json.Marshal(tasks)

	if err := redisClient.Client.Set(redisCtx, cacheKey, jsonData, 30*time.Second).Err(); err != nil {
		fmt.Println("Erro ao salvar no Redis:", err)
	}

	// -----------------------------------------------------------------------
	// 4. Atualiza o histórico (últimas 10 buscas)
	// -----------------------------------------------------------------------
	if err := pushSearchHistory(historyKey, query); err != nil {
		fmt.Println("Erro ao atualizar histórico:", err)
	}

	return tasks, nil
}

// ---------------------------------------------------------------------------
// Função auxiliar para armazenar histórico das últimas 10 buscas
// ---------------------------------------------------------------------------
func pushSearchHistory(historyKey, query string) error {
	pipe := redisClient.Client.TxPipeline()

	pipe.LPush(redisCtx, historyKey, query)
	pipe.LTrim(redisCtx, historyKey, 0, 9) // mantém apenas as últimas 10

	_, err := pipe.Exec(redisCtx)
	return err
}

func GetSearchHistory(userID uint) ([]string, error) {
	if redisClient == nil || redisClient.Client == nil {
		return []string{}, nil
	}

	historyKey := "tmpro:search:history:" + strconv.Itoa(int(userID))

	items, err := redisClient.Client.LRange(redisCtx, historyKey, 0, 9).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return items, nil
}
