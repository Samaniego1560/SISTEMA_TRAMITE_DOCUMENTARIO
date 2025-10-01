import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormRequirementComponent } from './form-requirement.component';

describe('FormRequirementComponent', () => {
  let component: FormRequirementComponent;
  let fixture: ComponentFixture<FormRequirementComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FormRequirementComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(FormRequirementComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
