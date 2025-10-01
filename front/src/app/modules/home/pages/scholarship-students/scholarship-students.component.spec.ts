import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ScholarshipStudentsComponent } from './scholarship-students.component';

describe('ScholarshipStudentsComponent', () => {
  let component: ScholarshipStudentsComponent;
  let fixture: ComponentFixture<ScholarshipStudentsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ScholarshipStudentsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ScholarshipStudentsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
