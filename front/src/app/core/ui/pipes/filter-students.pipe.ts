import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'filterStudents',
  standalone: true,
  pure: false
})
export class FilterStudentsPipe implements PipeTransform {

  transform(students: any[], searchText: string): any[] {
    if (!students) return [];
    if (!searchText) return students;

    searchText = searchText.toLowerCase();
    
    return students.filter(student =>
      student.student.code.toLowerCase().includes(searchText) ||
      student.student.full_name.toLowerCase().includes(searchText)
    );
  }

}
