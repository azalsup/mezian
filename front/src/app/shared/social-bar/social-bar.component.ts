import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';

@Component({
  selector: 'app-social-bar',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './social-bar.component.html',
  styleUrls: ['./social-bar.component.scss'],
})
export class SocialBarComponent {
  protected readonly lang = inject(LangService);
}
