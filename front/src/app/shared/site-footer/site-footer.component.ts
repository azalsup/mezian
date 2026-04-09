import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { LangService } from '../../core/services/lang.service';
import {
  faFacebook,
  faWhatsapp,
  faTiktok,
  faXTwitter,
  faInstagram,
  faLinkedin,
} from '@fortawesome/free-brands-svg-icons';

@Component({
  selector: 'app-site-footer',
  standalone: true,
  imports: [CommonModule, FontAwesomeModule],
  templateUrl: './site-footer.component.html',
  styleUrls: ['./site-footer.component.scss'],
})
export class SiteFooterComponent {
  protected readonly lang = inject(LangService);

  readonly faFacebook  = faFacebook;
  readonly faWhatsapp  = faWhatsapp;
  readonly faTiktok    = faTiktok;
  readonly faXTwitter  = faXTwitter;
  readonly faInstagram = faInstagram;
  readonly faLinkedin  = faLinkedin;
}
