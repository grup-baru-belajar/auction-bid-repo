#!/usr/bin/env python3
"""
Race Condition Simulator untuk Bidding System
Mensimulasikan concurrent bids dan detect race conditions
"""

import concurrent.futures
import json
import requests
import time
import sys
from datetime import datetime
from decimal import Decimal
from typing import List, Dict, Any

# Configuration
API_URL = "http://localhost:8080"
AUCTION_ID = 1
NUM_CONCURRENT_BIDS = 10
AUTH_TOKEN = "your_token_here"  # Replace with valid token

class BidSimulator:
    def __init__(self, api_url: str, auction_id: int, auth_token: str):
        self.api_url = api_url
        self.auction_id = auction_id
        self.auth_token = auth_token
        self.results: List[Dict[str, Any]] = []
        self.lock_time = None
    
    def place_bid(self, bidder_id: int, bid_price: float) -> Dict[str, Any]:
        """Place a single bid and record result"""
        headers = {
            "Authorization": f"Bearer {self.auth_token}",
            "Content-Type": "application/json"
        }
        
        payload = {
            "auctionID": self.auction_id,
            "bidPrice": bid_price
        }
        
        start_time = time.time()
        
        try:
            response = requests.post(
                f"{self.api_url}/bids",
                json=payload,
                headers=headers,
                timeout=5
            )
            
            elapsed_time = time.time() - start_time
            
            result = {
                "bidder_id": bidder_id,
                "bid_price": bid_price,
                "status_code": response.status_code,
                "timestamp": datetime.now().isoformat(),
                "elapsed_time": elapsed_time,
                "accepted": response.status_code == 201,
                "response": response.json() if response.headers.get('content-type') == 'application/json' else response.text
            }
            
            return result
        
        except Exception as e:
            return {
                "bidder_id": bidder_id,
                "bid_price": bid_price,
                "status_code": None,
                "timestamp": datetime.now().isoformat(),
                "elapsed_time": time.time() - start_time,
                "accepted": False,
                "error": str(e)
            }
    
    def simulate_sequential_bids(self) -> None:
        """Test 1: Sequential bids to establish baseline"""
        print("\n" + "="*60)
        print("TEST 1: SEQUENTIAL BIDS (Baseline)")
        print("="*60)
        
        for i in range(NUM_CONCURRENT_BIDS):
            bid_price = 100 + i
            result = self.place_bid(i, bid_price)
            self.results.append(result)
            
            status = "✅ ACCEPTED" if result["accepted"] else "❌ REJECTED"
            print(f"Bid #{i}: Price={bid_price} | {status} | Time: {result['elapsed_time']:.3f}s")
    
    def simulate_concurrent_bids(self) -> None:
        """Test 2: Concurrent bids - detect race conditions"""
        print("\n" + "="*60)
        print("TEST 2: CONCURRENT BIDS (Race Condition Test)")
        print("="*60)
        
        self.results = []
        
        with concurrent.futures.ThreadPoolExecutor(max_workers=NUM_CONCURRENT_BIDS) as executor:
            futures = []
            
            # Synchronize bidders to start at approximately the same time
            for i in range(NUM_CONCURRENT_BIDS):
                bid_price = 100 + i
                future = executor.submit(self.place_bid, i, bid_price)
                futures.append(future)
            
            # Collect results
            for future in concurrent.futures.as_completed(futures):
                result = future.result()
                self.results.append(result)
        
        # Sort by timestamp for analysis
        self.results.sort(key=lambda x: x["timestamp"])
        
        for result in self.results:
            status = "✅ ACCEPTED" if result["accepted"] else "❌ REJECTED"
            print(f"Bid #{result['bidder_id']}: Price={result['bid_price']} | {status} | Time: {result['elapsed_time']:.3f}s")
    
    def simulate_burst_bids(self, num_bursts: int = 3, delay_ms: int = 100) -> None:
        """Test 3: Burst of bids with controlled delays"""
        print("\n" + "="*60)
        print("TEST 3: BURST BIDS (Multiple waves)")
        print("="*60)
        
        self.results = []
        price = 100
        
        for burst in range(num_bursts):
            print(f"\n--- Burst {burst + 1} ---")
            burst_results = []
            
            with concurrent.futures.ThreadPoolExecutor(max_workers=NUM_CONCURRENT_BIDS // num_bursts) as executor:
                futures = []
                
                for i in range(NUM_CONCURRENT_BIDS // num_bursts):
                    bid_price = price + i
                    bidder_id = burst * (NUM_CONCURRENT_BIDS // num_bursts) + i
                    future = executor.submit(self.place_bid, bidder_id, bid_price)
                    futures.append(future)
                
                for future in concurrent.futures.as_completed(futures):
                    result = future.result()
                    self.results.append(result)
                    burst_results.append(result)
            
            for result in burst_results:
                status = "✅ ACCEPTED" if result["accepted"] else "❌ REJECTED"
                print(f"  Bid #{result['bidder_id']}: Price={result['bid_price']} | {status}")
            
            price += NUM_CONCURRENT_BIDS // num_bursts
            if burst < num_bursts - 1:
                print(f"Waiting {delay_ms}ms before next burst...")
                time.sleep(delay_ms / 1000)
    
    def analyze_results(self) -> None:
        """Analyze results for race conditions"""
        print("\n" + "="*60)
        print("ANALYSIS")
        print("="*60)
        
        if not self.results:
            print("No results to analyze")
            return
        
        accepted = [r for r in self.results if r["accepted"]]
        rejected = [r for r in self.results if not r["accepted"]]
        
        print(f"\nTotal Bids: {len(self.results)}")
        print(f"✅ Accepted: {len(accepted)} ({len(accepted)/len(self.results)*100:.1f}%)")
        print(f"❌ Rejected: {len(rejected)} ({len(rejected)/len(self.results)*100:.1f}%)")
        
        if accepted:
            print(f"\nAccepted Bids:")
            for bid in accepted:
                print(f"  - Bidder {bid['bidder_id']}: ${bid['bid_price']}")
        
        if len(self.results) > 1:
            times = sorted([r["elapsed_time"] for r in self.results])
            print(f"\nResponse Time Analysis:")
            print(f"  Min: {times[0]:.3f}s")
            print(f"  Avg: {sum(times)/len(times):.3f}s")
            print(f"  Max: {times[-1]:.3f}s")
        
        # Detect potential race conditions
        print(f"\nRace Condition Indicators:")
        
        # Check if multiple bids were accepted for same price
        prices = [r["bid_price"] for r in accepted]
        if len(prices) != len(set(prices)):
            print(f"  ⚠️  DETECTED: Multiple bids at same price level")
        
        # Check if bids with lower price were accepted after higher price
        for i in range(len(self.results) - 1):
            if (self.results[i]["accepted"] and self.results[i+1]["accepted"] and 
                self.results[i]["bid_price"] > self.results[i+1]["bid_price"]):
                print(f"  ⚠️  DETECTED: Lower bid accepted after higher bid")
                break
    
    def export_results(self, filename: str = "race_condition_results.json") -> None:
        """Export results to file"""
        with open(filename, 'w') as f:
            json.dump(self.results, f, indent=2, default=str)
        print(f"\n✅ Results exported to {filename}")


def main():
    print("🔥 Bidding System - Race Condition Simulator")
    print("="*60)
    print(f"API URL: {API_URL}")
    print(f"Auction ID: {AUCTION_ID}")
    print(f"Concurrent Bids: {NUM_CONCURRENT_BIDS}")
    
    simulator = BidSimulator(API_URL, AUCTION_ID, AUTH_TOKEN)
    
    try:
        # Run all tests
        simulator.simulate_sequential_bids()
        time.sleep(2)  # Wait before next test
        
        simulator.simulate_concurrent_bids()
        simulator.analyze_results()
        time.sleep(2)  # Wait before next test
        
        simulator.simulate_burst_bids(num_bursts=3, delay_ms=500)
        simulator.analyze_results()
        
        # Export results
        simulator.export_results()
        
    except KeyboardInterrupt:
        print("\n\n⚠️  Test interrupted by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
