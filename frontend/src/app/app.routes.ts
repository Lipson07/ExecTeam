import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { guestGuard } from './core/guards/guest.guard';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./features/meetpage/meetpage.component').then((m) => m.MeetpageComponent),
  },
  {
    path: 'auth',
    canMatch: [guestGuard],
    loadChildren: () => import('./features/auth/auth.routes').then((m) => m.authRoutes),
  },
  {
    path: 'app',
    canMatch: [authGuard],
    loadChildren: () => import('./features/homepage/homepage.routes').then((m) => m.homepageRoutes),
  },
  { path: '**', redirectTo: '' },
];
