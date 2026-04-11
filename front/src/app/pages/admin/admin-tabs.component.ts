import { Component, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-admin-tabs',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, CommonModule],
  templateUrl: './admin-tabs.component.html',
  styleUrls: ['./admin-tabs.component.scss'],
})
export class AdminTabsComponent {
  private readonly auth = inject(AuthService);
  protected readonly isStaff = this.auth.isStaff;
}