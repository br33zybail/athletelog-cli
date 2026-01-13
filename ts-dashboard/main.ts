import Chart, { ChartConfiguration } from 'chart.js/auto';

interface Workout {
    date: string;
    exercise: string;
    weight: number;
    reps: number;
    estimated_1rm?: number;
}

async function loadData() {
    try {
        const response = await fetch('/data/workouts.json');
        if (!response.ok) throw new Error('Failed to load workouts.json');
        const history: Workout[] = await response.json();

        // Populate table
        const tbody = document.getElementById('table-body')!;
        tbody.innerHTML = '';
        history.forEach(w => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${w.date}</td>
                <td>${w.exercise}</td>
                <td>${w.weight.toFixed(1)}</td>
                <td>${w.reps}</td>
                <td>${w.estimated_1rm ? w.estimated_1rm.toFixed(0) : '-'}</td>
            `;
            tbody.appendChild(row);
        });

        // Prepare datasets - one per exercise
        const exercises = [...new Set(history.map(w => w.exercise))];
        const datasets = exercises.map(ex => {
            const exData = history.filter(w => w.exercise === ex);
            return {
                label: ex.charAt(0).toUpperCase() + ex.slice(1), //Capitalize
                data: exData.map(w => w.weight),
                borderColor: getRandomColor(),
                backgroundColor: getRandomColor(0.2),
                tension: 0.3,
                fill: true,
                hidden: false
            };
        })

        // Helper to get random color
        function getRandomColor(opacity = 1) {
            const r = Math.floor(Math.random() * 255);
            const g = Math.floor(Math.random() * 255);
            const b = Math.floor(Math.random() * 255);
            return `rgba(${r}, ${g}, ${b}, ${opacity})`;
        }

        // Chart.js weight trend
        const canvas = document.getElementById('weightChart') as HTMLCanvasElement;
        const ctx = canvas.getContext('2d')!;
        
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: history.map(w => w.date).sort(),
                datasets: datasets
            },
            options: {
                responsive: true,
                plugins: {
                    legend: { position: 'top' },
                    title: { display: true, text: 'Weight Progress by Exercise' }
                },
                scales: {
                    y: {
                        beginAtZero: false,
                        title: { display: true, text: 'Weight (lbs)' }
                    },
                    x: {
                        title: { display: true, text: 'Date' }
                    }
                }
            }       
        });
    }   catch (err) {
        console.error(err);
        document.body.innerHTML += '<p style="color:red;">Error loading data. Make sure data/workouts.json exists.</p>';
    }
}

loadData();