import { Component, OnInit, OnDestroy } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { StockService } from '../../services/stock.service';
import { Holding } from '../../models/holdings.model';
import { Stock } from '../../models/stock.model';
import { Router } from '@angular/router';
import { CONSTANTS } from '../../config/constants';
import { Subject, takeUntil } from 'rxjs';

@Component({
  selector: 'app-portfolio',
  templateUrl: './portfolio.component.html',
  styleUrls: ['./portfolio.component.css']
})
export class PortfolioComponent implements OnInit, OnDestroy {
  username: string = '';
  selectedStock: Stock | null = null;
  openStates: Map<string, boolean> = new Map();
  holdings: Holding[] = [];
  errorFetchingHoldings: boolean = false;
  constants = CONSTANTS;
  private destroy$ = new Subject<void>();

  constructor(
    private authService: AuthService,
    private stockService: StockService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.fetchHoldings();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  fetchHoldings(): void {
    this.stockService.getHoldings().subscribe({
      next: (response) => {
        if (response && response.holdings && response.holdings.length > 0) {
          this.holdings = response.holdings;
          response.holdings.forEach(holding => {
            this.openStates.set(holding.stockID, false);
          });
          this.errorFetchingHoldings = false;
        } else {
          this.holdings = [];
          this.errorFetchingHoldings = false;
        }
      },
      error: (error) => {
        console.error('Error fetching holdings:', error);
        this.errorFetchingHoldings = true;
        this.holdings = [];
      }
    });
  }

  toggleStock(stockID: string): void {
    const isOpen = this.openStates.get(stockID) || false;
    this.openStates.set(stockID, !isOpen);

    if (!isOpen) {
      this.fetchStockDetails(stockID);
    } else {
      this.selectedStock = null;
    }
  }

  fetchStockDetails(stockID: string): void {
    this.stockService.getStockDetails(stockID)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (stockDetails) => {
          this.selectedStock = stockDetails;
        },
        error: (error) => {
          console.error('Error fetching stock details:', error);
          this.selectedStock = null;
        }
      });
  }

  onBuyStock(stockId: string): void {
    this.router.navigate(['/order'], { queryParams: { action: 'Buy', stockId: stockId } });
  }

  onSellStock(stockId: string): void {
    this.router.navigate(['/order'], { queryParams: { action: 'Sell', stockId: stockId } });
  }
}