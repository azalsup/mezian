import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavbarComponent } from './shared/navbar/navbar.component';
import { AuthModalComponent } from './shared/auth-modal/auth-modal.component';
import { LangService } from './core/services/lang.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, NavbarComponent, AuthModalComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  host: {
    '[dir]': 'lang.isRtl() ? "rtl" : "ltr"',
  },
})
export class AppComponent {
  readonly lang = inject(LangService);
}
