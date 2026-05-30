import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../../core/service/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
})
export class RegisterComponent {
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);

  readonly registerForm: FormGroup = this.fb.group({
    name: ['', [Validators.required, Validators.minLength(3)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    confirmPassword: ['', [Validators.required, Validators.minLength(6)]],
  });

  protected isSubmit = false;
  protected serverError = '';

  isValidated(controlName: string): boolean {
    const control = this.registerForm.get(controlName);
    return ((control?.touched || this.isSubmit) && control?.invalid) || false;
  }

  onSubmit(): void {
    this.isSubmit = true;
    this.serverError = '';

    if (this.registerForm.invalid) return;

    const { name, email, password, confirmPassword } = this.registerForm.value;

    if (password !== confirmPassword) {
      this.serverError = 'Пароли не совпадают';
      return;
    }

    this.authService.register(name, email, password).subscribe({
      next: () => {
        this.router.navigate(['/auth/verify-email'], { queryParams: { email } });
      },
      error: (err) => {
        this.serverError = err.error?.message || 'Ошибка регистрации';
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
