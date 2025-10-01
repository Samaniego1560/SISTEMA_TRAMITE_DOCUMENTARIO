import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GraphicsAreaComponent } from './graphics-area.component';

describe('GraphicsAreaComponent', () => {
  let component: GraphicsAreaComponent;
  let fixture: ComponentFixture<GraphicsAreaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [GraphicsAreaComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(GraphicsAreaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
