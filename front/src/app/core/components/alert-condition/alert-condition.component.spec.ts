import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AlertConditionComponent } from './alert-condition.component';

describe('AlertConditionComponent', () => {
  let component: AlertConditionComponent;
  let fixture: ComponentFixture<AlertConditionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AlertConditionComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AlertConditionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
