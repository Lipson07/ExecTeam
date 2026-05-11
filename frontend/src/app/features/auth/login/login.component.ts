import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
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
  loginForm: FormGroup;
  isSubmit = false;
  serverError = '';

  constructor(
    private router: Router,
    private fb: FormBuilder,
    private authService: AuthService,
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  get f() {
    return this.loginForm.controls;
  }

  isValidated(controlName: string): boolean {
    const control = this.loginForm.get(controlName);
    return ((control?.touched || this.isSubmit) && control?.invalid) || false;
  }

  onSubmit() {
    this.isSubmit = true;
    this.serverError = '';

    if (this.loginForm.invalid) {
      Object.keys(this.loginForm.controls).forEach((key) => {
        this.loginForm.get(key)?.markAsTouched();
      });
      return;
    }

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

  loginWithGitHub() {
    this.authService.loginWithGitHub();
  }

  loginWithGoogle() {
    this.authService.loginWithGoogle();
  }

  onBack() {
    this.router.navigate(['/']);
  }

  onReg() {
    this.router.navigate(['/auth/register']);
  }
}
