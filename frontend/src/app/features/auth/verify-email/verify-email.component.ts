import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../../../core/service/auth.service';

@Component({
  selector: 'app-verify-email',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './verify-email.component.html',
  styleUrls: ['./verify-email.component.scss'],
})
export class VerifyEmailComponent implements OnInit {
  verifyForm: FormGroup;
  isSubmit = false;
  serverError = '';
  successMessage = '';
  submittedEmail = '';

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService,
  ) {
    this.verifyForm = this.fb.group({
      code: ['', [Validators.required, Validators.minLength(6)]],
    });
  }

  ngOnInit(): void {
    this.submittedEmail = this.route.snapshot.queryParamMap.get('email') || '';
  }

  isValidated(): boolean {
    const control = this.verifyForm.get('code');
    return ((control?.touched || this.isSubmit) && control?.invalid) || false;
  }

  onSubmit() {
    this.isSubmit = true;
    this.serverError = '';
    this.successMessage = '';

    if (this.verifyForm.invalid) {
      this.verifyForm.get('code')?.markAsTouched();
      return;
    }

    const code = this.verifyForm.get('code')?.value;
    const userId = Number(this.route.snapshot.queryParamMap.get('userId'));

    this.authService.verifyEmail(this.submittedEmail, code, userId).subscribe({
      next: () => {
        this.successMessage = 'Почта подтверждена!';
        setTimeout(() => {
          this.router.navigate(['/app']);
        }, 1500);
      },
      error: (err) => {
        this.serverError = err.error?.message || 'Неверный код';
      },
    });
  }

  resendCode() {
    this.authService.resendCode(this.submittedEmail).subscribe({
      next: () => {
        this.successMessage = 'Код отправлен повторно';
      },
      error: (err) => {
        this.serverError = err.error?.message || 'Ошибка отправки';
      },
    });
  }

  onBack() {
    this.router.navigate(['/auth/login']);
  }
}
