import { NgClass, NgForOf } from '@angular/common';
import { Component, Output, EventEmitter, Input } from '@angular/core';


@Component({
    selector: 'app-pagination',
    standalone: true,
    imports: [NgClass, NgForOf],
    templateUrl: './pagination.html',
    styleUrl: './pagination.scss',
    providers: []
})
export class PaginationComponent {
    @Input() currentPage = 1;
    @Input() totalPages = 1;
    @Output() pageChange = new EventEmitter<number>();

    totalPagesArray: number[] = [];

    ngOnChanges(): void {
        this.totalPagesArray = Array.from({ length: this.totalPages }, (_, i) => i + 1);
    }

    goToPage(page: number | string) {
        if (typeof page === 'number' && page >= 1 && page <= this.totalPages) {
          this.pageChange.emit(page);
        }
      }

    getDisplayedPages(): (number | string)[] {
        const pages: (number | string)[] = [];
        const total = this.totalPages;
        const current = this.currentPage;
      
        if (total <= 4) {
          for (let i = 1; i <= total; i++) pages.push(i);
        } else {
          pages.push(1); // siempre mostrar primera
      
          if (current > 3) pages.push('...');
      
          const start = Math.max(2, current - 1);
          const end = Math.min(total - 1, current + 1);
      
          for (let i = start; i <= end; i++) pages.push(i);
      
          if (current < total - 2) pages.push('...');
      
          pages.push(total); // siempre mostrar última
        }
      
        return pages;
      }
      

}


