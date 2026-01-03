import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { finalize } from 'rxjs/operators';
import { ChangeDetectorRef, NgZone } from '@angular/core';

import {
  ReactiveFormsModule,
  FormBuilder,
  Validators,
  FormGroup,
} from '@angular/forms';
import {
  Subject,
  debounceTime,
  distinctUntilChanged,
  takeUntil,
} from 'rxjs';

import { TaskService, Task } from '../../core/tasks/task.service';
import { AuthService } from '../../core/auth/auth.service';
import { Router } from '@angular/router';

import { AiService } from '../../core/ai/ai.service';

type AiModalType = 'titles' | 'improve' | null;
type AiTarget = 'create' | 'edit';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './dashboard.html',
  styleUrls: ['./dashboard.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  tasks: Task[] = [];
  history: string[] = [];

  // ✅ loaders separados
  tasksLoading = false;   // só para LISTA (e só quando você quiser mostrar)
  actionLoading = false;  // criar/editar/deletar/buscar

  errorMsg = '';

  editingId: number | null = null;

  createForm!: FormGroup;
  editForm!: FormGroup;
  searchForm!: FormGroup;

  // =========================
  // IA (somente 2 features)
  // =========================
  aiLoading = false;
  aiError = '';
  aiModal: AiModalType = null;

  aiTitleOptions: string[] = [];
  aiImprovedText = '';
  aiImprovedBullets: string[] = [];

  private aiTarget: AiTarget = 'create';

  private destroy$ = new Subject<void>();

  constructor(
    private fb: FormBuilder,
    private taskService: TaskService,
    private auth: AuthService,
    private router: Router,
    private ai: AiService,
    private zone: NgZone,
    private cdr: ChangeDetectorRef
  ) {
    this.createForm = this.fb.group({
      title: ['', [Validators.required, Validators.minLength(2)]],
      description: [''],
      status: ['TODO', Validators.required],
      priority: ['MEDIUM', Validators.required],
    });

    this.editForm = this.fb.group({
      title: ['', [Validators.required, Validators.minLength(2)]],
      description: [''],
      status: ['TODO', Validators.required],
      priority: ['MEDIUM', Validators.required],
    });

    this.searchForm = this.fb.group({
      q: [''],
    });
  }

  ngOnInit(): void {
    // this.loadTasks(true); // se quiser carregar ao entrar e mostrar loader, descomenta
    this.loadHistory();

    this.searchForm
      .get('q')!
      .valueChanges.pipe(
        debounceTime(300),
        distinctUntilChanged(),
        takeUntil(this.destroy$)
      )
      .subscribe((q) => {
        const value = (q ?? '').trim();
        if (!value) {
          // ✅ quando limpa o campo (automático), recarrega SEM mostrar "Carregando..."
          this.loadTasks(false);
          return;
        }
        this.search(value);
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // =========================
  // IA helpers
  // =========================
  private resetAiState(): void {
    this.aiLoading = false;
    this.aiError = '';
    this.aiModal = null;
    this.aiTitleOptions = [];
    this.aiImprovedText = '';
    this.aiImprovedBullets = [];
  }

  private openAiModal(type: AiModalType, target: AiTarget): void {
    this.aiTarget = target;
    this.aiModal = type;
  }

  closeAiModal(): void {
    this.resetAiState();
  }

  private getFormByTarget(target: AiTarget): FormGroup {
    return target === 'create' ? this.createForm : this.editForm;
  }

  // =========================
  // LISTA
  // =========================
  /**
   * showSpinner = true  -> mostra "Carregando..." (uso do botão Atualizar)
   * showSpinner = false -> não mostra "Carregando..." (fluxos automáticos)
   */
  loadTasks(showSpinner: boolean = true): void {
    if (showSpinner) this.tasksLoading = true;
    this.errorMsg = '';

    this.taskService.getTasks().pipe(
      finalize(() => {
        // garante atualização visual
        this.zone.run(() => {
          this.tasksLoading = false;
          this.cdr.detectChanges();
        });
      })
    ).subscribe({
      next: (tasks) => {
        this.zone.run(() => {
          this.tasks = tasks || [];
          this.cdr.detectChanges();
        });
      },
      error: () => {
        this.zone.run(() => {
          this.errorMsg = 'Erro ao carregar tasks.';
          this.tasks = [];
          this.cdr.detectChanges();
        });
      },
    });
  }


  loadHistory(): void {
    this.taskService.getSearchHistory().subscribe({
      next: (items) => (this.history = items || []),
      error: () => (this.history = []),
    });
  }

  // =========================
  // CRUD
  // =========================
  onCreate(): void {
    if (this.createForm.invalid) return;

    this.actionLoading = true;
    this.errorMsg = '';

    const payload = this.createForm.value;

    this.taskService
      .createTask(payload)
      .pipe(finalize(() => (this.actionLoading = false)))
      .subscribe({
        next: () => {
          this.createForm.reset({
            title: '',
            description: '',
            status: 'TODO',
            priority: 'MEDIUM',
          });

          this.resetAiState();

          const q = (this.searchForm.get('q')?.value ?? '').trim();
          if (q) this.search(q);
          else this.loadTasks(false); // ✅ sem spinner
        },
        error: () => {
          this.errorMsg = 'Erro ao criar task.';
        },
      });
  }

  startEdit(t: Task): void {
    this.editingId = t.id;

    this.editForm.reset({
      title: t.title ?? '',
      description: t.description ?? '',
      status: t.status ?? 'TODO',
      priority: t.priority ?? 'MEDIUM',
    });

    this.resetAiState();
  }

  cancelEdit(): void {
    this.editingId = null;
    this.editForm.reset();
    this.resetAiState();
  }

  saveEdit(id: number): void {
    if (this.editForm.invalid) return;

    this.actionLoading = true;
    this.errorMsg = '';

    this.taskService
      .updateTask(id, this.editForm.value)
      .pipe(finalize(() => (this.actionLoading = false)))
      .subscribe({
        next: () => {
          this.editingId = null;
          this.resetAiState();

          const q = (this.searchForm.get('q')?.value ?? '').trim();
          if (q) this.search(q);
          else this.loadTasks(false); // ✅ sem spinner
        },
        error: () => {
          this.errorMsg = 'Erro ao atualizar task.';
        },
      });
  }

  deleteTask(id: number): void {
    this.actionLoading = true;
    this.errorMsg = '';

    this.taskService
      .deleteTask(id)
      .pipe(finalize(() => (this.actionLoading = false)))
      .subscribe({
        next: () => {
          const q = (this.searchForm.get('q')?.value ?? '').trim();
          if (q) this.search(q);
          else this.loadTasks(false); // ✅ sem spinner
        },
        error: () => {
          this.errorMsg = 'Erro ao deletar task.';
        },
      });
  }

  // =========================
  // SEARCH (Redis)
  // =========================
  search(q: string): void {
    this.actionLoading = true;
    this.errorMsg = '';

    this.taskService
      .searchTasks(q)
      .pipe(
        finalize(() => {
          this.zone.run(() => {
            this.actionLoading = false;
            this.cdr.detectChanges();
          });
        })
      )
      .subscribe({
        next: (tasks) => {
          this.zone.run(() => {
            this.tasks = tasks || [];
            this.loadHistory();
            this.cdr.detectChanges();
          });
        },
        error: () => {
          this.zone.run(() => {
            this.errorMsg = 'Erro ao buscar tasks.';
            this.cdr.detectChanges();
          });
        },
      });
  }

  clickHistory(item: string): void {
    this.searchForm.patchValue({ q: item }, { emitEvent: true });
  }

  clearSearch(): void {
    this.searchForm.patchValue({ q: '' }, { emitEvent: true });
    this.loadTasks(false); // ✅ sem spinner
    this.loadHistory();
  }

  // =========================
  // IA actions (2 features)
  // =========================
  onSuggestTitle(target: AiTarget): void {
    if (this.aiLoading) return;

    this.aiError = '';
    const form = this.getFormByTarget(target);

    const desc = (form.get('description')?.value ?? '').toString().trim();
    if (desc.length < 10) {
      this.aiError =
        'Escreva uma descrição (mín. 10 caracteres) para gerar títulos.';
      return;
    }

    this.aiLoading = true;
    this.ai.suggestTitle(desc).pipe(finalize(() => (this.aiLoading = false))).subscribe({
      next: (res) => {
        this.aiTitleOptions = res?.titles || [];
        this.openAiModal('titles', target);
      },
      error: () => {
        this.aiError = 'Falha ao sugerir títulos.';
      },
    });
  }

  applyTitle(title: string): void {
    const form = this.getFormByTarget(this.aiTarget);
    form.get('title')?.setValue(title);
    this.closeAiModal();
  }

  onImproveDescription(target: AiTarget): void {
    if (this.aiLoading) return;

    this.aiError = '';
    const form = this.getFormByTarget(target);

    const title = (form.get('title')?.value ?? '').toString().trim();
    const desc = (form.get('description')?.value ?? '').toString().trim();

    if (title.length < 3) {
      this.aiError =
        'Preencha o título (mín. 3 caracteres) para melhorar a descrição.';
      return;
    }
    if (desc.length < 10) {
      this.aiError =
        'Preencha a descrição (mín. 10 caracteres) para melhorar.';
      return;
    }

    this.aiLoading = true;
    this.ai.improveDescription(title, desc).pipe(finalize(() => (this.aiLoading = false))).subscribe({
      next: (res) => {
        this.aiImprovedText = res?.improved_description || '';
        this.aiImprovedBullets = res?.bullets || [];
        this.openAiModal('improve', target);
      },
      error: () => {
        this.aiError = 'Falha ao melhorar a descrição.';
      },
    });
  }

  applyImprovedDescription(): void {
    const form = this.getFormByTarget(this.aiTarget);
    const improved = (this.aiImprovedText ?? '').toString().trim();

    if (improved) {
      form.get('description')?.setValue(improved);
    }
    this.closeAiModal();
  }

  // =========================
  // LOGOUT
  // =========================
  logout(): void {
    this.auth.logout();
    this.router.navigate(['/login']);
  }
}
