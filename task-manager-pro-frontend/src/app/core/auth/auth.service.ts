import { Injectable, inject, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environments';
import { Observable, tap } from 'rxjs';
import { isPlatformBrowser } from '@angular/common';

type LoginResponse = { token: string };

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly tokenKey = 'tmpro_token';
  private readonly baseUrl = environment.apiUrl;

  private http = inject(HttpClient);
  private platformId = inject(PLATFORM_ID);

  private get isBrowser(): boolean {
    return isPlatformBrowser(this.platformId);
  }

  register(payload: { name: string; email: string; password: string }) {
  return this.http.post('/api/auth/register', payload);
}


  login(email: string, password: string): Observable<LoginResponse> {
  return this.http
    .post<LoginResponse>('/api/auth/login', { email, password })
    .pipe(tap((res) => this.setToken(res.token)));
}

  logout(): void {
    if (!this.isBrowser) return;
    localStorage.removeItem(this.tokenKey);
  }

  getToken(): string | null {
    if (!this.isBrowser) return null;
    return localStorage.getItem(this.tokenKey);
  }

  isLoggedIn(): boolean {
    return !!this.getToken();
  }

  private setToken(token: string): void {
    if (!this.isBrowser) return;
    localStorage.setItem(this.tokenKey, token);
  }
}
