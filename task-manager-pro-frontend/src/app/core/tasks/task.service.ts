import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environments';
import { Observable } from 'rxjs';

export type Task = {
  id: number;
  title: string;
  description: string;
  status: string;
  priority: string;
};

@Injectable({ providedIn: 'root' })
export class TaskService {
  private baseUrl = `${environment.apiUrl}/tasks`;

  constructor(private http: HttpClient) {}

  getTasks(): Observable<Task[]> {
    return this.http.get<Task[]>(this.baseUrl);
  }
}
