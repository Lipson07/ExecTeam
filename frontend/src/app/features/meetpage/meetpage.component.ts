import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-meetpage',
  imports: [],
  templateUrl: './meetpage.component.html',
  styleUrl: './meetpage.component.scss',
})
export class MeetpageComponent {
  constructor(private router: Router) {}
  onReg() {
    this.router.navigate(['/auth/register']);
  }
  onLogin() {
    this.router.navigate(['/auth/login']);
  }
}
