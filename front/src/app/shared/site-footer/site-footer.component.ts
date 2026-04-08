import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SocialBarComponent } from '../social-bar/social-bar.component';
import { LangService } from '../../core/services/lang.service';

@Component({
  selector: 'app-site-footer',
  standalone: true,
  imports: [CommonModule, SocialBarComponent],
  templateUrl: './site-footer.component.html',
  styleUrls: ['./site-footer.component.scss'],
})
export class SiteFooterComponent {
  protected readonly lang = inject(LangService);
}
