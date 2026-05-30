import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-meetpage',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './meetpage.component.html',
  styleUrls: ['./meetpage.component.scss'],
})
export class MeetpageComponent {}
