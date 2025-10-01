import { ComponentFixture, TestBed } from '@angular/core/testing';

import { HouseholdTargetingComponent } from './household-targeting.component';

describe('HouseholdTargetingComponent', () => {
  let component: HouseholdTargetingComponent;
  let fixture: ComponentFixture<HouseholdTargetingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [HouseholdTargetingComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(HouseholdTargetingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
