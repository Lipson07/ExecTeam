import { Component, Output, EventEmitter, signal } from '@angular/core';

@Component({
  selector: 'app-project-create',
  standalone: true,
  templateUrl: './project-create.component.html',
  styleUrls: ['./project-create.component.scss'],
})
export class ProjectCreateComponent {
  @Output() readonly close = new EventEmitter<void>();

  readonly currentStep = signal(1);
  readonly totalSteps = 4;

  readonly repoAction = signal<'existing' | 'new' | 'skip' | null>(null);
  readonly selectedRepoId = signal<number | null>(null);
  readonly selectedOwnerId = signal<number | null>(1);
  readonly newRepoVisibility = signal<'private' | 'public'>('private');

  protected readonly repos = signal([
    { id: 1, name: 'execTeam/backend', desc: 'Основной API на Go' },
    { id: 2, name: 'execTeam/frontend', desc: 'Angular-приложение' },
    { id: 3, name: 'execTeam/bot', desc: 'Telegram-бот' },
  ]);

  protected readonly owners = signal([
    { id: 1, name: 'Алексей Иванов', email: 'alexey@execteam.ru', initials: 'АИ' },
    { id: 2, name: 'Мария Петрова', email: 'maria@execteam.ru', initials: 'МП' },
    { id: 3, name: 'Дмитрий Сидоров', email: 'dmitry@execteam.ru', initials: 'ДС' },
  ]);

  setRepoAction(action: 'existing' | 'new' | 'skip'): void {
    this.repoAction.set(action);
    if (action !== 'existing') {
      this.selectedRepoId.set(null);
    }
  }

  selectRepo(id: number): void {
    this.selectedRepoId.set(id);
  }

  selectOwner(id: number): void {
    this.selectedOwnerId.set(id);
  }

  setRepoVisibility(visibility: 'private' | 'public'): void {
    this.newRepoVisibility.set(visibility);
  }

  closeModal(): void {
    this.close.emit();
  }

  nextStep(): void {
    if (this.currentStep() === 2) {
      if (this.repoAction() === 'skip') {
        this.currentStep.set(4);
        return;
      }
      if (!this.repoAction()) return;
    }
    if (this.currentStep() < this.totalSteps) {
      this.currentStep.update((s) => s + 1);
    }
  }

  prevStep(): void {
    if (this.currentStep() === 4 && this.repoAction() === 'skip') {
      this.currentStep.set(2);
      return;
    }
    if (this.currentStep() > 1) {
      this.currentStep.update((s) => s - 1);
    }
  }

  onSubmit(): void {
    this.close.emit();
  }
}
