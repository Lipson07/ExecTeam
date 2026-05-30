import { Component, signal } from '@angular/core';
import { HeaderComponent } from '../../shared/components/header/header.component';
import { FooterComponent } from '../../shared/components/footer/footer.component';
import { CommonModule } from '@angular/common';
import { ProjectCreateComponent } from '../../shared/modals/project-create/project-create.component';

@Component({
  selector: 'app-homepage',
  standalone: true,
  imports: [CommonModule, HeaderComponent, FooterComponent, ProjectCreateComponent],
  templateUrl: './homepage.component.html',
  styleUrl: './homepage.component.scss',
})
export class HomepageComponent {
  readonly showCreateModal = signal(false);

  openCreateModal(): void {
    this.showCreateModal.set(true);
  }

  closeCreateModal(): void {
    this.showCreateModal.set(false);
  }
}
