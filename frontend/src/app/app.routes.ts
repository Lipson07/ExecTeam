import { Routes } from '@angular/router';
import { MeetpageComponent } from './features/meetpage/meetpage.component';
import { AuthGuard } from './core/guards/auth.guard';
import { GuestGuard } from './core/guards/guest.guard';

export const routes: Routes = [
  { path: '', component: MeetpageComponent },
  {
    path: 'auth',
    canActivate: [GuestGuard],
    loadChildren: () => import('./features/auth/auth.routes').then((m) => m.authRoutes),
  },
  {
    path: 'app',
    canActivate: [AuthGuard],
    loadChildren: () => import('./features/homepage/homepage.routes').then((m) => m.homepageRoutes),
  },
  { path: '**', redirectTo: '' },
];
