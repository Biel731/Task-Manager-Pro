import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';

import { TaskService, Task } from '../../core/tasks/task.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.html',
})
export class DashboardComponent {
  loading = true;
  tasks: Task[] = [];
  error = '';

  constructor(private taskService: TaskService) {
    this.loadTasks();
  }

  loadTasks(): void {
    this.taskService.getTasks().subscribe({
      next: (tasks) => {
        this.tasks = tasks;
        this.loading = false;
      },
      error: () => {
        this.error = 'Erro ao carregar tasks';
        this.loading = false;
      },
    });
  }
}
