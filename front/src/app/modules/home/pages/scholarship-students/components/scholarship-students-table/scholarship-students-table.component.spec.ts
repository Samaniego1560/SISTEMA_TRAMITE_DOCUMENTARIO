import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ScholarshipStudentsTableComponent } from './scholarship-students-table.component';

describe('ScholarshipStudentsTableComponent', () => {
  let component: ScholarshipStudentsTableComponent;
  let fixture: ComponentFixture<ScholarshipStudentsTableComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ScholarshipStudentsTableComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ScholarshipStudentsTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
