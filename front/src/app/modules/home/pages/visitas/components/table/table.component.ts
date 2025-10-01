// table.component.ts
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

export interface TableColumn {
  header: string;
  field: string;
  sortable?: boolean;
  filterable?: boolean;
}

export interface TableAction {
  label: string;
  icon: string;
  action: 'view' | 'edit' | 'delete' | 'register';
}

export interface FilterOption {
  value: any;
  label: string;
}

export interface TableFilter {
  field: string;
  label: string;
  options: FilterOption[];
  dependsOn?: string; // Campo del que depende este filtro (para filtros anidados)
  childFilters?: { [key: string]: FilterOption[] }; // Opciones de filtro hijo basadas en el valor del padre
  customFilter?: (data: any[], filterValue: string) => any[]; // Función de filtrado personalizada
}

@Component({
  selector: 'app-table',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
})
export class TableComponent {
  @Input() columns: TableColumn[] = [];
  @Input() data: any[] = [];
  @Input() actions: TableAction[] = [
    { label: 'Ver', icon: 'eye', action: 'view' },
    { label: 'Editar', icon: 'pencil', action: 'edit' },
    { label: 'Eliminar', icon: 'trash', action: 'delete' }
  ];
 
  @Input() filters: TableFilter[] = [];
  activeFilters: { [key: string]: any } = {};
  
  @Output() viewItem = new EventEmitter<any>();
  @Output() editItem = new EventEmitter<any>();
  @Output() deleteItem = new EventEmitter<any>();
  @Output() registerItem = new EventEmitter<any>();
  
  @Input() itemsPerPage = 10;
  currentPage = 1;
  searchTerm = '';
  filteredData: any[] = [];
  Math = Math;
  
  ngOnInit() {
    this.filteredData = [...this.data];
  }
  
  ngOnChanges() {
    this.filteredData = [...this.data];
    this.applyFilters();
  }

  onActionClick(action: string, item: any) {
    switch(action) {
      case 'view':
        this.viewItem.emit(item);
        break;
      case 'edit':
        this.editItem.emit(item);
        break;
      case 'delete':
        this.deleteItem.emit(item);
        break;
      case 'register':
        this.registerItem.emit(item);
        break;
    }
  }
  
  get totalPages(): number {
    return Math.ceil(this.filteredData.length / this.itemsPerPage);
  }
  
  get paginatedData(): any[] {
    const startIndex = (this.currentPage - 1) * this.itemsPerPage;
    return this.filteredData.slice(startIndex, startIndex + this.itemsPerPage);
  }
  
  onPageChange(page: number) {
    if (page >= 1 && page <= this.totalPages) {
      this.currentPage = page;
    }
  }
  
  applyFilters() {
    let result = [...this.data];
    
    if (this.searchTerm.trim()) {
      const term = this.searchTerm.toLowerCase();
      result = result.filter(item => {
        return Object.keys(item).some(key => {
          const value = item[key];
          if (value !== null && value !== undefined) {
            return value.toString().toLowerCase().includes(term);
          }
          return false;
        });
      });
    }

    Object.keys(this.activeFilters).forEach(field => {
      const filterValue = this.activeFilters[field];
      if (filterValue !== null && filterValue !== undefined && filterValue !== '') {
        // Buscar si hay un filtro personalizado para este campo
        const filterConfig = this.filters.find(f => f.field === field);
        
        if (filterConfig && filterConfig.customFilter) {
          // Usar el filtro personalizado
          result = filterConfig.customFilter(result, filterValue);
        } else {
          // Usar el filtro estándar
          result = result.filter(item => item[field] === filterValue);
        }
      }
    });
    
    this.filteredData = result;
    this.currentPage = 1; 
  }

  onFilterChange(field: string, value: any) {
    this.activeFilters[field] = value;
    this.filters.forEach(filter => {
      if (filter.dependsOn === field) {
        this.activeFilters[filter.field] = '';
      }
    });
    
    this.applyFilters();
  }

  getFilterOptions(filter: TableFilter): FilterOption[] {
    if (!filter.dependsOn) {
      return filter.options;
    }
    if (filter.childFilters) {
      const parentValue = this.activeFilters[filter.dependsOn];
      if (parentValue && parentValue !== '') {
        return filter.childFilters[parentValue] || [];
      } else {
        if (filter.options && filter.options.length > 0) {
          return filter.options;
        } else {
          const allOptionsMap = new Map<string, FilterOption>();
          Object.values(filter.childFilters).forEach(options => {
            options.forEach(option => {
              allOptionsMap.set(option.value, option);
            });
          });
          return Array.from(allOptionsMap.values());
        }
      }
    }
    
    return filter.options;
  }
  
  isFilterDisabled(filter: TableFilter): boolean {
    return false;
  }
  getVisiblePages(): number[] {
  const pages: number[] = [];
  const delta = 2; 
  
  const start = Math.max(1, this.currentPage - delta);
  const end = Math.min(this.totalPages, this.currentPage + delta);
  
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }
  
  return pages;
}
}