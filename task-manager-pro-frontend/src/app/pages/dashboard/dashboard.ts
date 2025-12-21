import { Component, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
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

  loading = false;
  errorMsg = '';

  editingId: number | null = null;

  createForm!: FormGroup;
  editForm!: FormGroup;
  searchForm!: FormGroup;

  private destroy$ = new Subject<void>();

  constructor(
    private fb: FormBuilder,
    private taskService: TaskService,
    private auth: AuthService,
    private router: Router
  ) {
    // cria forms AQUI pra não dar "fb used before initialization"
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
    this.loadTasks();
    this.loadHistory();

    // escuta busca (Redis) com debounce
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
          this.loadTasks();
          return;
        }
        this.search(value);
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // -------------------------
  // LOADERS
  // -------------------------
  loadTasks(): void {
    this.loading = true;
    this.errorMsg = '';

    this.taskService.getTasks().subscribe({
      next: (tasks) => {
        this.tasks = tasks;
        this.loading = false;
      },
      error: () => {
        this.errorMsg = 'Erro ao carregar tasks.';
        this.loading = false;
      },
    });
  }

loadHistory(): void {
  this.taskService.getSearchHistory().subscribe({
    next: (items) => (this.history = items),
    error: () => (this.history = []),
  });
}

  // -------------------------
  // CRUD
  // -------------------------
  onCreate(): void {
    if (this.createForm.invalid) return;

    this.loading = true;
    this.errorMsg = '';

    const payload = this.createForm.value;

    this.taskService.createTask(payload).subscribe({
      next: () => {
        this.createForm.reset({
          title: '',
          description: '',
          status: 'TODO',
          priority: 'MEDIUM',
        });

        // se tiver busca ativa, mantém o contexto
        const q = (this.searchForm.get('q')?.value ?? '').trim();
        if (q) this.search(q);
        else this.loadTasks();

        this.loading = false;
      },
      error: () => {
        this.errorMsg = 'Erro ao criar task.';
        this.loading = false;
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
  }

  cancelEdit(): void {
    this.editingId = null;
    this.editForm.reset();
  }

  saveEdit(id: number): void {
    if (this.editForm.invalid) return;

    this.loading = true;
    this.errorMsg = '';

    this.taskService.updateTask(id, this.editForm.value).subscribe({
      next: () => {
        this.editingId = null;

        const q = (this.searchForm.get('q')?.value ?? '').trim();
        if (q) this.search(q);
        else this.loadTasks();

        this.loading = false;
      },
      error: () => {
        this.errorMsg = 'Erro ao atualizar task.';
        this.loading = false;
      },
    });
  }

  deleteTask(id: number): void {
    this.loading = true;
    this.errorMsg = '';

    this.taskService.deleteTask(id).subscribe({
      next: () => {
        const q = (this.searchForm.get('q')?.value ?? '').trim();
        if (q) this.search(q);
        else this.loadTasks();

        this.loading = false;
      },
      error: () => {
        this.errorMsg = 'Erro ao deletar task.';
        this.loading = false;
      },
    });
  }

  // -------------------------
  // SEARCH (Redis)
  // -------------------------
  search(q: string): void {
    this.loading = true;
    this.errorMsg = '';

    this.taskService.searchTasks(q).subscribe({
      next: (tasks) => {
        this.tasks = tasks;
        this.loading = false;
        this.loadHistory();
      },
      error: () => {
        this.errorMsg = 'Erro ao buscar tasks.';
        this.loading = false;
      },
    });
  }

  clickHistory(item: string): void {
    this.searchForm.patchValue({ q: item }, { emitEvent: true });
  }

  clearSearch(): void {
    this.searchForm.patchValue({ q: '' }, { emitEvent: true });
    this.loadTasks();
    this.loadHistory();
  }

  // -------------------------
  // LOGOUT
  // -------------------------
  logout(): void {
    this.auth.logout();
    this.router.navigate(['/login']);
  }
}
