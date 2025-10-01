import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ScholarshipStudentsFormComponent } from './scholarship-students-form.component';

describe('ScholarshipStudentsFormComponent', () => {
  let component: ScholarshipStudentsFormComponent;
  let fixture: ComponentFixture<ScholarshipStudentsFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ScholarshipStudentsFormComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ScholarshipStudentsFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
