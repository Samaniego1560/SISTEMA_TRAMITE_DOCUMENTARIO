import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GraphicsStatisticComponent } from './graphics-statistic.component';

describe('GraphicsStatisticComponent', () => {
  let component: GraphicsStatisticComponent;
  let fixture: ComponentFixture<GraphicsStatisticComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GraphicsStatisticComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(GraphicsStatisticComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
