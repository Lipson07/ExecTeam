import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/service/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);

  readonly loginForm: FormGroup = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
  });

  protected isSubmit = false;
  protected serverError = '';

  isValidated(controlName: string): boolean {
    const control = this.loginForm.get(controlName);
    return ((control?.touched || this.isSubmit) && control?.invalid) || false;
  }

  onSubmit(): void {
    this.isSubmit = true;
    this.serverError = '';

    if (this.loginForm.invalid) return;

    const { email, password } = this.loginForm.value;

    this.authService.login(email, password).subscribe({
      next: (res) => {
        if (res.need_verification) {
          this.router.navigate(['/auth/verify-email'], {
            queryParams: { email, userId: res.user_id },
          });
        } else {
          this.router.navigate(['/app']);
        }
      },
      error: (err) => {
        this.serverError = err.error?.message || 'Неверный email или пароль';
      },
    });
  }

  loginWithGitHub(): void {
    this.authService.loginWithGitHub();
  }

  loginWithGoogle(): void {
    this.authService.loginWithGoogle();
  }
}
