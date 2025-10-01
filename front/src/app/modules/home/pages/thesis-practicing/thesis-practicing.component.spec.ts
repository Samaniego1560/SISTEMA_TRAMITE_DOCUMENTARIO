import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ThesisPracticingComponent } from './thesis-practicing.component';

describe('ThesisPracticingComponent', () => {
  let component: ThesisPracticingComponent;
  let fixture: ComponentFixture<ThesisPracticingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ThesisPracticingComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ThesisPracticingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
