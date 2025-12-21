import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from '../../../environments/environments';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export type Tag = {
  id: number;
  user_id: number;
  name: string;
};

export type Task = {
  id: number;
  user_id: number;
  title: string;
  description: string;
  status: 'TODO' | 'DOING' | 'DONE' | string;
  priority: 'LOW' | 'MEDIUM' | 'HIGH' | string;
  due_date?: string | null;
  tags?: Tag[];
  created_at?: string;
  updated_at?: string;
};

export type CreateTaskPayload = {
  title: string;
  description: string;
  status?: string;   // default no backend se você quiser
  priority?: string; // default no backend se você quiser
  due_date?: string | null;
  tags?: string[];   // backend espera names (CreateTaskInput.Tags)
};

export type UpdateTaskPayload = {
  title?: string;
  description?: string;
  status?: string;
  priority?: string;
  due_date?: string | null;
  tags?: string[]; // se enviar, substitui as tags
};

export type TaskListFilters = {
  status?: string;
  priority?: string;
  tags?: string; // no teu repo.go está "tags" (singular) como query param
  q?: string;    // busca em title/description (ListTasks usa filter.Query)
};

@Injectable({ providedIn: 'root' })
export class TaskService {
  private readonly baseUrl = `${environment.apiUrl}/tasks`;

  constructor(private http: HttpClient) {}

  // -------------------------
  // CRUD
  // -------------------------

  getTasks(filters?: TaskListFilters): Observable<Task[]> {
    let params = new HttpParams();

    if (filters?.status) params = params.set('status', filters.status);
    if (filters?.priority) params = params.set('priority', filters.priority);
    if (filters?.tags) params = params.set('tags', filters.tags);
    if (filters?.q) params = params.set('q', filters.q);

    return this.http.get<Task[]>(this.baseUrl, { params });
  }

  getTaskById(id: number): Observable<Task> {
    return this.http.get<Task>(`${this.baseUrl}/${id}`);
  }

  createTask(payload: CreateTaskPayload): Observable<Task> {
    // garante arrays limpos
    const body: CreateTaskPayload = {
      ...payload,
      tags: payload.tags?.map((t) => t.trim()).filter(Boolean),
    };

    return this.http.post<Task>(this.baseUrl, body);
  }

  updateTask(id: number, payload: UpdateTaskPayload): Observable<Task> {
    const body: UpdateTaskPayload = {
      ...payload,
      tags: payload.tags?.map((t) => t.trim()).filter(Boolean),
    };

    return this.http.put<Task>(`${this.baseUrl}/${id}`, body);
  }

  deleteTask(id: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/${id}`);
  }

  // -------------------------
  // Redis Search endpoint
  // -------------------------
  searchTasks(q: string): Observable<Task[]> {
    const params = new HttpParams().set('q', q);
    return this.http.get<Task[]>(`${this.baseUrl}/search`, { params });
  }

  getSearchHistory(): Observable<string[]> {
  return this.http
    .get<{ history: string[] }>(`${this.baseUrl}/search/history`)
    .pipe(map((res) => res?.history ?? []));
}


}
