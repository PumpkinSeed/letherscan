export const prerender = true;

export async function load() {
    try {
        const response = await fetch('http://localhost:8080/blocks');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        return {
            blocks: data.blocks
        };
    } catch (error) {
        console.error('Error fetching blocks:', error);
        return {
            blocks: []
        };
    }
} 