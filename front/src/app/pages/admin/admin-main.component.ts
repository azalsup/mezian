import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { AdminTabsComponent } from './admin-tabs.component';

@Component({
  selector: 'app-admin-main',
  standalone: true,
  imports: [RouterOutlet, AdminTabsComponent],
  templateUrl: './admin-main.component.html',
  styleUrls: ['./admin-main.component.scss'],
})
export class AdminMainComponent {}