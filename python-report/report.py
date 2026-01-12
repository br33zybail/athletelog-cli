import sys
import json
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
from pathlib import Path

def main ():
    if len(sys.argv) != 2:
        print("Usage: python report.py <path_to_workouts.json>")
        sys.exit(1)
        
    json_path = Path(sys.argv[1])
    if not json_path.exists():
        print(f"File not found: {json_path}")
        sys.exit(1)
        
    with open(json_path, "r") as f:
        data = json.load(f)
        
    if not data:
        print("No workouts yet.")
        return
    
    df = pd.DataFrame(data)
    df['exercise'] = df['exercise'].str.lower().str.strip()
    df['date'] = pd.to_datetime(df['date'])
    df = df.sort_values('date')
    
    # Calculate volume (weight * reps)
    df['volume'] = df['weight'] * df['reps']
    
    # Group by exercise for plotting
    exercise = df['exercise'].unique()
    
    # Plot weight over time
    plt.figure(figsize=(10, 6))
    sns.lineplot(data=df, x='date', y='weight', hue='exercise', marker='o')
    plt.title('Weight Progress Over Time')
    plt.xlabel('Date')
    plt.ylabel('Weight')
    plt.grid(True)
    plt.xticks(rotation=45)
    
    # Save plot
    
    output_dir = json_path.parent
    plot_path = output_dir / 'report.png'
    plt.savefig(plot_path, bbox_inches='tight')
    plt.close()
    
    # Text summary
    
    summary = []
    for ex in exercise:
        ex_df = df[df['exercise'] == ex]
        summary.append(f"{ex.capitalize()}:")
        summary.append(f"  - Entries: {len(ex_df)}")
        summary.append(f"  - Max Weight: {ex_df['weight'].max():.1f} lbs")
        summary.append(f"  - Total Volume: {ex_df['volume'].sum():.1f} lbs")
        summary.append("")
        
    total_volume = df['volume'].sum()
    summary.append(f"Overall total volume: {total_volume:.1f} lbs")
    
    summary_path = output_dir / 'summary.txt'
    with open(summary_path, 'w') as f:
        f.write("\n".join(summary))
        
    print(f"Report generated: {plot_path}")
    print(f"Summary generated: {summary_path}")
    
if __name__ == "__main__":
    main()