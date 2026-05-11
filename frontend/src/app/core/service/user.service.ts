import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Observable } from 'rxjs';
interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  isVerified: boolean;
}
@Injectable({
  providedIn: 'root',
})
export class UserService {
  constructor(private api: ApiService) {}
  getMe(): Observable<User> {
    return this.api.get<User>('/users/me');
  }
  updateMe(data: { name?: string; email: string; avatar?: string }): Observable<User> {
    return this.api.put<User>('/users/me', data);
  }
  getAll(): Observable<User[]> {
    return this.api.get<User[]>('/users');
  }

  searchByName(query: string): Observable<User[]> {
    return this.api.get<User[]>('/users/search', { q: query });
  }

  getById(id: string): Observable<User> {
    return this.api.get<User>(`/users/${id}`);
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/users/${id}`);
  }
}
