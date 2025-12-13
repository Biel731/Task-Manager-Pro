import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environments';
import { Observable, tap } from 'rxjs';

type LoginResponse = {
  token: string;
};

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly tokenKey = 'tmpro_token';
  private readonly baseUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  register(email: string, password: string): Observable<unknown> {
    return this.http.post(`${this.baseUrl}/auth/register`, { email, password });
  }

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http
      .post<LoginResponse>(`${this.baseUrl}/auth/login`, { email, password })
      .pipe(tap((res) => this.setToken(res.token)));
  }

  logout(): void {
    localStorage.removeItem(this.tokenKey);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }

  private setToken(token: string): void {
    localStorage.setItem(this.tokenKey, token);
  }
}
