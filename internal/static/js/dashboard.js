  function app() {
            return {
                currentView: 'dashboard',
                filter: 'all', 
                tasks: [],
                selectedTask: null,
                isLoading: true,
                showModal: false,
                newTask: { name: '', payload: '', runAt: '' },
                nowString: '', 
                userHasEditedTime: false,
                toast: { show: false, message: '', type: 'success' },
                connected: false,
                socket: null,

                init() {
                    this.fetchTasks();
                    this.connectWS();
                    this.updateTimeState();
                    setInterval(() => this.updateTimeState(), 1000);
                },

                connectWS() {
                    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
                    this.socket = new WebSocket(`${proto}//${window.location.host}/ws`);

                    this.socket.onopen = () => { this.connected = true; };
                    
                    this.socket.onmessage = (event) => {
                        const data = JSON.parse(event.data);
                        this.handleEvent(data);
                    };

                    this.socket.onclose = () => {
                        this.connected = false;
                        setTimeout(() => this.connectWS(), 3000);
                    };
                },

                handleEvent(data) {
                    const normalizeId = (obj) => String(obj.ID || obj.id || '');

                    if (data.type === 'TASK_CREATED') {
                        const exists = this.tasks.some(t => normalizeId(t) === normalizeId(data.payload));
                        if (!exists) {
                            this.tasks.unshift(data.payload);
                        }
                    } 
                    else if (data.type === 'TASK_UPDATED') {
                        const targetId = normalizeId(data.payload);
                        
                        const idx = this.tasks.findIndex(t => normalizeId(t) === targetId);
                        
                        if (idx !== -1) {
                            this.tasks[idx] = { ...data.payload }; 
                        }

                        if (this.selectedTask && normalizeId(this.selectedTask) === targetId) {
                            this.selectedTask = { ...data.payload };
                        }
                    }
                },

                async fetchTasks(silent = false) {
                    if (!silent) this.isLoading = true;
                    try {
                        const res = await fetch('/api/tasks');
                        if (res.status === 401) window.location.href = '/login';
                        const data = await res.json();
                        this.tasks = data.tasks || [];
                    } catch (e) { console.error(e); } 
                    finally { if (!silent) this.isLoading = false; }
                },

                async fetchSingleTask(id, silent = false) {
                    try {
                        const res = await fetch(`/api/tasks/${id}`);
                        const data = await res.json();
                        if (data.task) this.selectedTask = JSON.parse(JSON.stringify(data.task));
                    } catch (e) { console.error(e); }
                },

                async submitTask() {
                    const selected = new Date(this.newTask.runAt).getTime();
                    const now = new Date().getTime();
                    if (selected < now - 30000) {
                        this.showToast("Cannot schedule task in the past", "error");
                        return;
                    }
                    let runAtValue = this.userHasEditedTime ? new Date(this.newTask.runAt).toISOString() : null;

                    try {
                        let pl = this.newTask.payload;
                        try { if(pl) pl = JSON.parse(pl); } catch(e){}
                        
                        const res = await fetch('/api/tasks', {
                            method: 'POST',
                            headers: {'Content-Type': 'application/json'},
                            body: JSON.stringify({
                                name: this.newTask.name,
                                payload: pl,
                                run_at: runAtValue
                            })
                        });

                        if (res.ok) {
                            this.showToast("Task Created", "success");
                            this.closeModal();
                            
                        } else throw new Error();
                    } catch(e) { this.showToast("Failed to create", "error"); }
                },

                loadDetail(id) {
                    this.selectedTask = null;
                    this.currentView = 'detail';
                    this.fetchSingleTask(id);
                },
                goBack() {
                    this.currentView = 'dashboard';
                    this.selectedTask = null;
                },
                async performAction(action) {
                    if (!confirm(`Are you sure you want to ${action} this task?`)) return;
                    const id = this.selectedTask.ID || this.selectedTask.id;
                    try {
                        const res = await fetch(`/api/tasks/${id}/${action}`, { method: 'POST' });
                        if (res.ok) {
                            this.showToast(`Task ${action} successful`, 'success');
                        } else throw new Error();
                    } catch (e) { this.showToast(`Failed to ${action}`, 'error'); }
                },
                openModal() {
                    this.userHasEditedTime = false;
                    this.updateTimeState(); 
                    this.newTask.name = ''; 
                    this.newTask.payload = '';
                    this.showModal = true;
                },
                closeModal() { this.showModal = false; },
                
                get filteredTasks() {
                    if (this.filter === 'all') return this.tasks;
                    return this.tasks.filter(t => {
                        const s = t.State || t.state;
                        if (this.filter === 'running') return s === 'leased';
                        if (this.filter === 'dead') return s === 'dead';
                        if (this.filter === 'pending') return s === 'pending'|| s === 'retry';
                        return s === this.filter;
                    });
                },
                get filterLabel() { return this.filter === 'all' ? 'All' : this.filter.charAt(0).toUpperCase() + this.filter.slice(1); },
                get stats() {
                    return {
                        total: this.tasks.length,
                        pending: this.tasks.filter(t => ['pending', 'retry'].includes(t.State||t.state)).length,
                        running: this.tasks.filter(t => (t.State||t.state) === 'leased').length,
                        completed: this.tasks.filter(t => (t.State||t.state) === 'completed').length,
                        dead: this.tasks.filter(t => (t.State||t.state) === 'dead').length
                    }
                },
                getISOStringLocal() {
                    const now = new Date();
                    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
                    return now.toISOString().slice(0, 19);
                },
                updateTimeState() {
                    const iso = this.getISOStringLocal();
                    this.nowString = iso;
                    if (this.showModal && !this.userHasEditedTime) this.newTask.runAt = iso;
                },
                getRetryCount(t) { return (t.Retries !== undefined ? t.Retries : t.retries) || 0; },
                getMaxRetries(t) { return (t.MaxRetries !== undefined ? t.MaxRetries : t.max_retries) || 3; },
                canRetry(t) { 
                    if(!t) return false; 
                    const s = (t.State || t.state);
                    const r = this.getRetryCount(t);
                    const max = this.getMaxRetries(t);
                    return ['failed', 'dead', 'canceled', 'retry'].includes(s) && r < max; 
                },
                canCancel(t) { if(!t) return false; return ['pending', 'leased', 'retry'].includes(t.State||t.state); },
                formatDate(d) { if (!d || d.startsWith('0001')) return 'Never'; return new Date(d).toLocaleString(); },
                formatPayload(p) { return typeof p === 'object' ? JSON.stringify(p, null, 2) : (p || 'No Payload'); },
                showToast(msg, type) {
                    this.toast = { show: true, message: msg, type: type };
                    setTimeout(() => this.toast.show = false, 3000);
                }
            }
        }