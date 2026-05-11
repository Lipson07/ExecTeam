import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../../core/service/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
})
export class RegisterComponent {
  registerForm: FormGroup;
  isSubmit = false;
  serverError = '';

  constructor(
    private authService: AuthService,
    private router: Router,
    private fb: FormBuilder,
  ) {
    this.registerForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      confirmPassword: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  isValidated(controlName: string): boolean {
    const control = this.registerForm.get(controlName);
    return ((control?.touched || this.isSubmit) && control?.invalid) || false;
  }

  get f() {
    return this.registerForm.controls;
  }

  passwordMatch(): boolean {
    const password = this.registerForm.get('password')?.value;
    const confirmPassword = this.registerForm.get('confirmPassword')?.value;
    return password === confirmPassword;
  }

  onSubmit() {
    this.isSubmit = true;
    this.serverError = '';

    if (this.registerForm.invalid) {
      Object.keys(this.registerForm.controls).forEach((key) => {
        this.registerForm.get(key)?.markAsTouched();
      });
      return;
    }

    if (!this.passwordMatch()) {
      this.serverError = 'Пароли не совпадают';
      return;
    }

    const { name, email, password } = this.registerForm.value;

    this.authService.register(name, email, password).subscribe({
      next: () => {
        this.router.navigate(['/auth/verify-email'], {
          queryParams: { email },
        });
      },
      error: (err) => {
        this.serverError = err.error?.message || 'Ошибка регистрации';
      },
    });
  }

  loginWithGitHub() {
    this.authService.loginWithGitHub();
  }

  loginWithGoogle() {
    this.authService.loginWithGoogle();
  }

  onBack() {
    this.router.navigate(['/']);
  }

  onLogin() {
    this.router.navigate(['/auth/login']);
  }
}
