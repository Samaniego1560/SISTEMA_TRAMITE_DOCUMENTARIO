import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DataStatisticSrqComponent } from './data-statistic-srq.component';

describe('DataStatisticSrqComponent', () => {
  let component: DataStatisticSrqComponent;
  let fixture: ComponentFixture<DataStatisticSrqComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DataStatisticSrqComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DataStatisticSrqComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
