import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StudentsResidencesComponent } from './students-residences.component';

describe('StudentsResidencesComponent', () => {
  let component: StudentsResidencesComponent;
  let fixture: ComponentFixture<StudentsResidencesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [StudentsResidencesComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(StudentsResidencesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
