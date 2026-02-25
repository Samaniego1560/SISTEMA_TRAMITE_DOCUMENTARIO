import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { StudentParentSurveyService } from '../student-parent-survey/student-parent-survey.service';

interface TotalsStats {
  total: number;
  has_children: number;
  no_children: number;
  gestation: number;
  pregnant: number;
  avg_children?: number | null;
}

interface SchoolStats {
  school: string;
  total: number;
  has_children: number;
  no_children: number;
  gestation: number;
  avg_children?: number | null;
}

interface ChildrenCountStat {
  children_count: number | null;
  total: number;
}

@Component({
  selector: 'app-student-parent-survey-admin',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './student-parent-survey-admin.component.html',
  styleUrl: './student-parent-survey-admin.component.scss'
})
export class StudentParentSurveyAdminComponent implements OnInit {
  loading = true;
  errorMessage = '';
  totals: TotalsStats | null = null;
  bySchool: SchoolStats[] = [];
  childrenCounts: ChildrenCountStat[] = [];

  constructor(private survey: StudentParentSurveyService) {}

  ngOnInit(): void {
    this.fetchStats();
  }

  fetchStats(): void {
    this.loading = true;
    this.errorMessage = '';
    this.survey.stats().subscribe({
      next: (data) => {
        this.totals = data?.totals ?? null;
        this.bySchool = data?.by_school ?? [];
        this.childrenCounts = data?.children_counts ?? [];
        this.loading = false;
      },
      error: () => {
        this.loading = false;
        this.errorMessage = 'No se pudieron cargar las estadísticas.';
      }
    });
  }

  downloadExcel(): void {
    this.survey.exportExcel().subscribe((blob) => {
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'encuesta_estudiantes.xlsx';
      link.click();
      window.URL.revokeObjectURL(url);
    });
  }

  percent(part: number, total: number): string {
    if (!total) {
      return '0%';
    }
    return `${((part / total) * 100).toFixed(1)}%`;
  }

  avgChildren(value?: number | null): string {
    if (value === null || value === undefined) {
      return '-';
    }
    return value.toFixed(2);
  }
}
