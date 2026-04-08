import { Component, Input, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LangService } from '../../core/services/lang.service';

@Component({
  selector: 'app-categories-bar',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './categories-bar.component.html',
  styleUrls: ['./categories-bar.component.scss'],
})
export class CategoriesBarComponent {
  protected readonly lang = inject(LangService);
  @Input() categories: Array<{ icon: string; slug: string; labelFr: string; labelAr: string }> = [];
}
