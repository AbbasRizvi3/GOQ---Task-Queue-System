function metrics() {
    return {
        stats: { 
            cpus: 0, 
            memory_alloc: 0, 
            memory_sys: 1, 
            goroutines: 0, 
            num_gc: 0, 
            uptime: "..." 
        },
        init() {
            this.fetchData();
            setInterval(() => this.fetchData(), 1000);
        },
        async fetchData() {
            try {
                const res = await fetch('/metrics?format=json');
                if (res.ok) {
                    this.stats = await res.json();
                }
            } catch(e) { 
                console.error("Metrics poll failed:", e); 
            }
        }
    }
}