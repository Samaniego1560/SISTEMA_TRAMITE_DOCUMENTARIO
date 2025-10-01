import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GraphicsPieComponent } from './graphics-pie.component';

describe('GraphicsPieComponent', () => {
  let component: GraphicsPieComponent;
  let fixture: ComponentFixture<GraphicsPieComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GraphicsPieComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(GraphicsPieComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
