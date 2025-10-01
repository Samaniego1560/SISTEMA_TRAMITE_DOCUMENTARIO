import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AlertDebtsComponent } from './alert-debts.component';

describe('AlertDebtsComponent', () => {
  let component: AlertDebtsComponent;
  let fixture: ComponentFixture<AlertDebtsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AlertDebtsComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AlertDebtsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
