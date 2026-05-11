import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, tap } from 'rxjs';

interface User {
  id: number;
  name: string;
  email: string;
  email_verified: boolean;
}

interface AuthResponse {
  token?: string;
  user?: User;
  message?: string;
}

interface LoginResponse {
  need_verification: boolean;
  user_id?: number;
  token?: string;
  user?: User;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private userSubject = new BehaviorSubject<User | null>(null);
  user$ = this.userSubject.asObservable();

  constructor(
    private api: ApiService,
    private router: Router,
  ) {
    const token = localStorage.getItem('token');
    if (token) {
      this.loadUser();
    }
  }

  loginWithGoogle(): void {
    window.location.href = 'http://localhost:8080/api/auth/google';
  }

  loginWithGitHub(): void {
    window.location.href = 'http://localhost:8080/api/auth/github';
  }

  handleOAuthCallback(token: string): void {
    localStorage.setItem('token', token);
    this.loadUser();
    this.router.navigate(['/app']);
  }

  register(name: string, email: string, password: string): Observable<AuthResponse> {
    return this.api.post<AuthResponse>('/register', { name, email, password });
  }

  login(email: string, password: string): Observable<LoginResponse> {
    return this.api.post<LoginResponse>('/login', { email, password }).pipe(
      tap((res) => {
        if (!res.need_verification && res.token) {
          localStorage.setItem('token', res.token);
          this.userSubject.next(res.user!);
        }
      }),
    );
  }

  verifyEmail(email: string, code: string, userId: number): Observable<AuthResponse> {
    return this.api.post<AuthResponse>('/verify-email', { email, code }).pipe(
      tap((res) => {
        if (res.token) {
          localStorage.setItem('token', res.token);
          this.userSubject.next(res.user!);
        }
      }),
    );
  }

  resendCode(email: string): Observable<{ message: string }> {
    return this.api.post<{ message: string }>('/resend-code', { email });
  }

  private loadUser(): void {
    this.api.get<User>('/users/me').subscribe({
      next: (user) => this.userSubject.next(user),
      error: () => this.logout(),
    });
  }

  logout(): void {
    localStorage.removeItem('token');
    this.userSubject.next(null);
    this.router.navigate(['/']);
  }

  isLoggedIn(): boolean {
    return !!localStorage.getItem('token');
  }

  getUser(): User | null {
    return this.userSubject.value;
  }
}
