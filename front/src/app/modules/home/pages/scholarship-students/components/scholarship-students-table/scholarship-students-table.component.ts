import {Component, effect, input, OnInit} from '@angular/core';
import {ScholarshipStudentsModel} from "../../models/scholarship-students.model";
import {FormsModule} from "@angular/forms";

@Component({
  selector: 'app-scholarship-students-table',
  standalone: true,
  imports: [
    FormsModule
  ],
  templateUrl: './scholarship-students-table.component.html',
  styleUrl: './scholarship-students-table.component.scss'
})
export class ScholarshipStudentsTableComponent {
  scholarshipStudents = input.required<ScholarshipStudentsModel[]>();
  protected leftLimit: number = 0;
  protected currentPage: number = 1;
  protected totalPage: number = 0;
  protected rightLimit: number = 10;
  public displayedStudents: ScholarshipStudentsModel[] = [];

  constructor() {
    effect(() => {
      this.paginate(this.scholarshipStudents())
      this.totalPage = this.scholarshipStudents().length;
    });
  }

  protected previousPage(): void {
    this.currentPage--;
    this.leftLimit = this.leftLimit - 10;
    this.rightLimit = this.leftLimit + 10;
    this.paginate();
  }

  protected paginate(data?: ScholarshipStudentsModel[]): void {
    const requests: ScholarshipStudentsModel[] = JSON.parse(JSON.stringify(data || this.scholarshipStudents()));
    this.displayedStudents = requests.slice(this.leftLimit, this.rightLimit);
  }

  protected nextPage(): void {
    this.leftLimit = this.currentPage * 10;
    this.rightLimit = this.leftLimit + 10;
    this.currentPage++;
    this.paginate();
  }

  findRequest($event: Event) {

  }
}
