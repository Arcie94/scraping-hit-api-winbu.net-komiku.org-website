import React from 'react';
import { Code } from 'lucide-react';
import EndpointDebugger from '../components/EndpointDebugger';

const Home = () => {
    const endpoints = [
        {
            category: "Komiku (Manga)",
            color: "from-yellow-500 to-orange-600",
            badge: "MANGA",
            items: [
                {
                    method: "GET",
                    path: "/api/v1/komiku/home",
                    description: "Latest & Popular manga",
                    params: []
                },
                {
                    method: "GET",
                    path: "/api/v1/komiku/search",
                    description: "Cari manga berdasarkan kata kunci",
                    params: [
                        { name: 'q', in: 'query', placeholder: 'Contoh: naruto', required: true }
                    ]
                },
                {
                    method: "GET",
                    path: "/api/v1/komiku/manga/{slug}",
                    description: "Detail manga + list chapter",
                    params: [
                        { name: 'slug', in: 'path', placeholder: 'Contoh: one-piece', required: true }
                    ]
                },
                {
                    method: "GET",
                    path: "/api/v1/komiku/chapter/{slug}",
                    description: "Gambar-gambar dari satu chapter",
                    params: [
                        { name: 'slug', in: 'path', placeholder: 'Contoh: one-piece-chapter-1', required: true }
                    ]
                },
                {
                    method: "GET",
                    path: "/api/v1/komiku/genres",
                    description: "List semua genre manga",
                    params: []
                },
            ]
        },
        {
            category: "Winbu (Anime)",
            color: "from-blue-500 to-purple-600",
            badge: "ANIME",
            items: [
                {
                    method: "GET",
                    path: "/api/v1/winbu/home",
                    description: "Top Series, Latest Anime, Movies",
                    params: []
                },
                {
                    method: "GET",
                    path: "/api/v1/winbu/search",
                    description: "Cari anime/drama berdasarkan kata kunci",
                    params: [
                        { name: 'q', in: 'query', placeholder: 'Contoh: jujutsu', required: true }
                    ]
                },
                {
                    method: "GET",
                    path: "/api/v1/winbu/detail/{slug}",
                    description: "Detail anime + list episode",
                    params: [
                        { name: 'slug', in: 'path', placeholder: 'Contoh: jujutsu-kaisen-season-2', required: true }
                    ]
                },
                {
                    method: "GET",
                    path: "/api/v1/winbu/episode/{slug}",
                    description: "Stream URLs + navigasi episode",
                    params: [
                        { name: 'slug', in: 'path', placeholder: 'Contoh: jujutsu-kaisen-s2-episode-1', required: true }
                    ]
                },
            ]
        }
    ];

    return (
        <div className="min-h-screen bg-gradient-to-br from-background via-gray-900 to-background">
            {/* Hero */}
            <div className="container mx-auto px-4 py-16 text-center">
                <div className="inline-block mb-4 px-4 py-2 bg-primary/10 border border-primary/30 rounded-full">
                    <Code className="inline w-5 h-5 text-primary mr-2" />
                    <span className="text-primary font-mono text-sm">API Explorer</span>
                </div>
                <h1 className="text-5xl md:text-7xl font-bold mb-4 bg-gradient-to-r from-white via-blue-100 to-purple-200 bg-clip-text text-transparent">
                    Komiku & Winbu API
                </h1>
                <p className="text-gray-400 text-lg max-w-2xl mx-auto mb-8">
                    Interactive REST API documentation for Manga & Anime data scraping
                </p>
                <div className="flex flex-wrap gap-4 justify-center text-sm">
                    <span className="px-4 py-2 bg-secondary/50 border border-white/10 rounded-lg">
                        🌐 Multi-language support
                    </span>
                    <span className="px-4 py-2 bg-secondary/50 border border-white/10 rounded-lg">
                        ⚡ Fast caching
                    </span>
                    <span className="px-4 py-2 bg-secondary/50 border border-white/10 rounded-lg">
                        📡 Real-time scraping
                    </span>
                </div>
            </div>

            {/* Endpoints */}
            <div className="container mx-auto px-4 pb-16 space-y-12">
                {endpoints.map((section, sIdx) => (
                    <section key={sIdx}>
                        <div className="mb-6 flex items-center gap-4">
                            <h2 className={`text-3xl font-bold bg-gradient-to-r ${section.color} bg-clip-text text-transparent`}>
                                {section.category}
                            </h2>
                            <span className={`px-3 py-1 text-xs font-bold rounded-full bg-gradient-to-r ${section.color} text-white`}>
                                {section.badge}
                            </span>
                        </div>

                        <div className="grid gap-4">
                            {section.items.map((endpoint, eIdx) => (
                                <EndpointDebugger
                                    key={eIdx}
                                    method={endpoint.method}
                                    path={endpoint.path}
                                    description={endpoint.description}
                                    params={endpoint.params}
                                />
                            ))}
                        </div>
                    </section>
                ))}
            </div>

            {/* Footer Info */}
            <div className="container mx-auto px-4 pb-12 space-y-6">
                {/* Usage Notes */}
                <div className="bg-yellow-500/10 border border-yellow-500/30 rounded-xl p-6">
                    <h3 className="text-yellow-400 font-bold mb-3 flex items-center gap-2">
                        ⚠️ Catatan Penggunaan
                    </h3>
                    <ul className="text-gray-300 text-sm space-y-2 list-disc list-inside">
                        <li>Klik pada endpoint untuk expand form debugger</li>
                        <li>Isi parameter yang diperlukan (bertanda *)</li>
                        <li>Klik "Kirim Request" untuk test API</li>
                        <li>Response akan ditampilkan dengan syntax highlighting</li>
                    </ul>
                </div>

                {/* Rate Limiting */}
                <div className="bg-blue-500/10 border border-blue-500/30 rounded-xl p-6">
                    <h3 className="text-blue-400 font-bold mb-3 flex items-center gap-2">
                        🛡️ Rate Limiting
                    </h3>
                    <ul className="text-gray-300 text-sm space-y-2 list-disc list-inside">
                        <li><strong>Limit:</strong> 60 requests per menit per IP address</li>
                        <li><strong>Response saat limit:</strong> HTTP 429 "Rate limit exceeded"</li>
                        <li><strong>Reset otomatis:</strong> Counter reset setiap menit</li>
                        <li><strong>Cache:</strong> Data di-cache selama 30 menit untuk performa optimal</li>
                    </ul>
                </div>

                {/* Tips & Troubleshooting */}
                <div className="bg-purple-500/10 border border-purple-500/30 rounded-xl p-6">
                    <h3 className="text-purple-400 font-bold mb-3 flex items-center gap-2">
                        💡 Tips & Troubleshooting
                    </h3>
                    <ul className="text-gray-300 text-sm space-y-2 list-disc list-inside">
                        <li><strong>Response null:</strong> Coba query berbeda atau tunggu cache expire (30 menit)</li>
                        <li><strong>Slow response:</strong> Request pertama ~2-3 detik (scraping), request selanjutnya ~100ms (cache)</li>
                        <li><strong>CORS error:</strong> API allow all origins untuk development, update untuk production</li>
                        <li><strong>Rate limited:</strong> Tunggu 60 detik atau gunakan IP berbeda untuk testing</li>
                        <li><strong>Slug format:</strong> Gunakan format lowercase dengan dash (contoh: <code>one-piece</code>)</li>
                    </ul>
                </div>
            </div>
        </div>
    );
};

export default Home;
