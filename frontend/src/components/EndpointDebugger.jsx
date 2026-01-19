import React, { useState } from 'react';
import { Send, Copy, Check, Loader, ChevronDown, ChevronUp } from 'lucide-react';
import axios from 'axios';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';

const EndpointDebugger = ({ method, path, description, params = [] }) => {
    const [isExpanded, setIsExpanded] = useState(false);
    const [inputValues, setInputValues] = useState({});
    const [response, setResponse] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [statusCode, setStatusCode] = useState(null);
    const [responseTime, setResponseTime] = useState(null);
    const [copied, setCopied] = useState(false);

    // Build URL dengan parameter
    const buildUrl = () => {
        let url = path;

        // Replace path params
        params.forEach(param => {
            if (param.in === 'path') {
                const value = inputValues[param.name] || `{${param.name}}`;
                url = url.replace(`{${param.name}}`, value);
            }
        });

        // Add query params
        const queryParams = params.filter(p => p.in === 'query' && inputValues[p.name]);
        if (queryParams.length > 0) {
            const queryString = queryParams
                .map(p => `${p.name}=${encodeURIComponent(inputValues[p.name])}`)
                .join('&');
            url += `?${queryString}`;
        }

        return url;
    };

    const handleInputChange = (paramName, value) => {
        setInputValues(prev => ({ ...prev, [paramName]: value }));
    };

    const handleSend = async () => {
        const path = buildUrl();
        const fullUrl = `http://localhost:3000${path}`;  // Hit port 3000 directly
        console.log('Sending request to:', fullUrl);

        setLoading(true);
        setError(null);
        setResponse(null);
        setStatusCode(null);
        setResponseTime(null);

        const startTime = performance.now();

        try {
            const res = await fetch(fullUrl, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                },
            });
            const endTime = performance.now();

            console.log('Response status:', res.status);
            console.log('Response headers:', res.headers);
            setStatusCode(res.status);
            setResponseTime(Math.round(endTime - startTime));

            // Get raw text first
            const rawText = await res.text();
            console.log('Raw response text:', rawText);
            console.log('Raw response length:', rawText.length);

            // Try to parse JSON
            if (rawText && rawText.trim()) {
                const data = JSON.parse(rawText);
                console.log('Parsed JSON data:', data);
                setResponse(data);
            } else {
                console.warn('Empty response body');
                setResponse(null);
            }

        } catch (err) {
            const endTime = performance.now();
            console.error('Request error:', err);

            setError(err.message || 'Network error');
            setStatusCode(0);
            setResponseTime(Math.round(endTime - startTime));
        } finally {
            setLoading(false);
        }
    };

    const copyResponse = () => {
        navigator.clipboard.writeText(JSON.stringify(response, null, 2));
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const previewUrl = `http://localhost${buildUrl()}`;

    return (
        <div className="group bg-secondary/30 backdrop-blur-sm border border-white/10 rounded-xl overflow-hidden hover:border-primary/50 transition-all duration-300">
            {/* Header - Always Visible */}
            <div
                className="p-6 cursor-pointer"
                onClick={() => setIsExpanded(!isExpanded)}
            >
                <div className="flex flex-col md:flex-row md:items-center gap-4">
                    {/* Method Badge */}
                    <span className="px-3 py-1 bg-green-500/20 text-green-400 border border-green-500/50 rounded font-mono text-sm font-bold w-fit">
                        {method}
                    </span>

                    {/* Path */}
                    <code className="flex-1 text-white font-mono text-sm md:text-base break-all">
                        {path}
                    </code>

                    {/* Expand Icon */}
                    <div className="flex items-center gap-2">
                        {isExpanded ? (
                            <ChevronUp className="w-5 h-5 text-primary" />
                        ) : (
                            <ChevronDown className="w-5 h-5 text-gray-400" />
                        )}
                    </div>
                </div>

                {/* Description */}
                <p className="text-gray-400 text-sm mt-3 pl-0 md:pl-20">
                    {description}
                </p>
            </div>

            {/* Debugger Panel - Expandable */}
            {isExpanded && (
                <div className="border-t border-white/10 p-6 bg-black/20 space-y-4">

                    {/* Input Parameters */}
                    {params.length > 0 && (
                        <div className="space-y-3">
                            <h4 className="text-sm font-bold text-gray-300">Parameter:</h4>
                            {params.map((param, idx) => (
                                <div key={idx}>
                                    <label className="text-xs text-gray-500 font-mono flex items-center gap-2 mb-1">
                                        {param.name}
                                        {param.required && <span className="text-red-400">*</span>}
                                    </label>
                                    <input
                                        type="text"
                                        placeholder={param.placeholder}
                                        value={inputValues[param.name] || ''}
                                        onChange={(e) => handleInputChange(param.name, e.target.value)}
                                        className="w-full bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-white text-sm focus:outline-none focus:border-primary transition-colors"
                                    />
                                </div>
                            ))}
                        </div>
                    )}

                    {/* URL Preview */}
                    <div>
                        <h4 className="text-xs text-gray-500 mb-2">Preview URL:</h4>
                        <code className="block bg-gray-900 border border-gray-700 rounded px-3 py-2 text-sm text-cyan-400 break-all">
                            {previewUrl}
                        </code>
                    </div>

                    {/* Send Button */}
                    <button
                        onClick={handleSend}
                        disabled={loading}
                        className="w-full md:w-auto px-6 py-3 bg-primary hover:bg-blue-600 disabled:bg-gray-600 text-white rounded-lg font-medium flex items-center justify-center gap-2 transition-all shadow-[0_0_15px_rgba(59,130,246,0.5)]"
                    >
                        {loading ? (
                            <>
                                <Loader className="w-5 h-5 animate-spin" /> Mengirim...
                            </>
                        ) : (
                            <>
                                <Send className="w-5 h-5" /> Kirim Request
                            </>
                        )}
                    </button>

                    {/* Response Section */}
                    {statusCode !== null && (
                        <div className="space-y-3">
                            {/* Status Bar */}
                            <div className="flex items-center gap-4 text-sm">
                                <span className={`px-3 py-1 rounded font-bold ${statusCode >= 200 && statusCode < 300
                                    ? 'bg-green-500/20 text-green-400 border border-green-500/50'
                                    : 'bg-red-500/20 text-red-400 border border-red-500/50'
                                    }`}>
                                    {statusCode >= 200 && statusCode < 300 ? '✅' : '❌'} {statusCode}
                                </span>
                                <span className="px-3 py-1 bg-yellow-500/20 text-yellow-400 border border-yellow-500/50 rounded font-mono">
                                    ⏱️ {responseTime}ms
                                </span>
                                {response && (
                                    <button
                                        onClick={copyResponse}
                                        className="ml-auto px-3 py-1 bg-gray-700 hover:bg-gray-600 text-white rounded flex items-center gap-2 transition-colors"
                                    >
                                        {copied ? (
                                            <>
                                                <Check className="w-4 h-4" /> Tersalin
                                            </>
                                        ) : (
                                            <>
                                                <Copy className="w-4 h-4" /> Salin
                                            </>
                                        )}
                                    </button>
                                )}
                            </div>

                            {/* Response Body */}
                            <div className="rounded-lg overflow-hidden border border-gray-700">
                                <div className="bg-gray-800 px-3 py-2 text-xs text-gray-400 font-mono">
                                    Response:
                                </div>
                                {error && !response ? (
                                    <div className="bg-red-900/20 p-4 text-red-400 text-sm">
                                        ❌ Error: {error}
                                    </div>
                                ) : (
                                    <SyntaxHighlighter
                                        language="json"
                                        style={vscDarkPlus}
                                        customStyle={{
                                            margin: 0,
                                            padding: '1rem',
                                            fontSize: '0.75rem',
                                            maxHeight: '400px',
                                            overflowY: 'auto',
                                        }}
                                    >
                                        {JSON.stringify(response, null, 2)}
                                    </SyntaxHighlighter>
                                )}
                            </div>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default EndpointDebugger;
