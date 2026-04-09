import { Component, Input, Output, EventEmitter, inject } from '@angular/core';
import { LangService } from '../../core/services/lang.service';

@Component({
  selector: 'app-ads-search-bar',
  standalone: true,
  templateUrl: './ads-search-bar.component.html',
})
export class AdsSearchBarComponent {
  readonly lang = inject(LangService);

  @Input() set query(v: string) { this._text = v; }
  @Output() readonly search = new EventEmitter<string>();

  _text = '';

  submit(): void { this.search.emit(this._text.trim()); }
}
