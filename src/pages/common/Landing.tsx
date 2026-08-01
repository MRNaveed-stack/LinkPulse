import { Link } from 'react-router-dom';
import { Link2, BarChart3, Sparkles, Zap, ArrowRight, Shield, Activity } from 'lucide-react';

export default function Landing() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50/50 via-slate-50 to-purple-50/50 text-gray-900 font-sans overflow-x-hidden">
      {/* Glow Effects */}
      <div className="absolute top-0 left-1/4 w-[500px] h-[500px] bg-indigo-500/5 rounded-full blur-3xl pointer-events-none"></div>
      <div className="absolute top-1/3 right-1/4 w-[600px] h-[600px] bg-purple-500/5 rounded-full blur-3xl pointer-events-none"></div>

      {/* Navbar */}
      <nav className="relative max-w-7xl mx-auto px-6 py-5 flex items-center justify-between border-b border-gray-200/80 backdrop-blur-md bg-white/60 z-10">
        <div className="flex items-center gap-2.5">
          <div className="p-2 bg-gradient-to-tr from-indigo-500 to-purple-600 rounded-xl shadow-lg shadow-indigo-500/30">
            <Activity className="h-6 w-6 text-white" />
          </div>
          <span className="text-xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-indigo-600 to-purple-600">
            LinkPulse
          </span>
        </div>

        <div className="flex items-center gap-4">
          <Link to="/login">
            <button className="border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 hover:text-gray-950 px-5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 cursor-pointer shadow-sm">
              Log In
            </button>
          </Link>
          <Link to="/register">
            <button className="bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 shadow-lg shadow-indigo-500/10 px-5 py-2.5 rounded-xl text-sm font-semibold text-white transition-all duration-200 cursor-pointer border border-transparent">
              Sign Up Free
            </button>
          </Link>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative max-w-7xl mx-auto px-6 pt-20 pb-24 text-center z-10">
        <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-indigo-50 border border-indigo-100 text-indigo-600 text-xs font-semibold uppercase tracking-wider mb-8">
          <Sparkles className="h-4 w-4" /> Next-Gen Link-In-Bio Solution
        </div>

        <h1 className="text-5xl md:text-7xl font-black tracking-tight leading-tight max-w-4xl mx-auto mb-6 bg-clip-text text-transparent bg-gradient-to-r from-slate-900 via-indigo-950 to-purple-950">
          Pulse Your Presence. <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600">
            One Link, Infinite Connections.
          </span>
        </h1>

        <p className="text-gray-600 text-lg md:text-xl leading-relaxed max-w-2xl mx-auto mb-10">
          The ultimate hub for creators, builders, and brands. Share everything you do in a single link, customize your visual identity, and track audience analytics in real-time.
        </p>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-20">
          <Link to="/register" className="w-full sm:w-auto">
            <button className="w-full sm:w-auto flex items-center justify-center gap-2 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 shadow-xl shadow-indigo-500/15 px-8 py-4 rounded-xl text-base font-bold text-white transition-all duration-200 group cursor-pointer border border-transparent">
              Claim Your Link
              <ArrowRight className="h-5 w-5 group-hover:translate-x-1 transition-transform" />
            </button>
          </Link>
          <Link to="/login" className="w-full sm:w-auto">
            <button className="w-full sm:w-auto border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 hover:text-gray-950 px-8 py-4 rounded-xl text-base font-bold transition-all duration-200 cursor-pointer shadow-sm">
              Explore Dashboard
            </button>
          </Link>
        </div>

        {/* Demo Visual Showcase */}
        <div className="relative max-w-4xl mx-auto bg-white/80 border border-gray-200/80 rounded-2xl p-4 shadow-2xl backdrop-blur-xl">
          <div className="absolute -inset-0.5 bg-gradient-to-r from-indigo-500/5 to-purple-500/5 rounded-2xl blur opacity-30"></div>
          <div className="relative bg-gray-50/50 rounded-xl overflow-hidden border border-gray-100 p-6 flex flex-col md:flex-row items-center justify-between gap-8">
            <div className="text-left max-w-md">
              <h3 className="text-2xl font-bold mb-3 bg-clip-text text-transparent bg-gradient-to-r from-slate-900 to-gray-800">
                A Beautiful Public Profile
              </h3>
              <p className="text-gray-500 text-sm leading-relaxed mb-6">
                Deliver a stunning experience for your mobile and desktop visitors. Your custom links, public bio, and active branding update in seconds.
              </p>
              <div className="flex flex-col gap-2">
                <div className="flex items-center gap-2 text-indigo-700 text-sm font-semibold">
                  <Zap className="h-4 w-4 text-indigo-500" /> Blazing fast Go redirections
                </div>
                <div className="flex items-center gap-2 text-indigo-700 text-sm font-semibold">
                  <Shield className="h-4 w-4 text-indigo-500" /> Built-in rate limiting and security
                </div>
              </div>
            </div>

            {/* Profile Preview Block */}
            <div className="w-full max-w-[280px] bg-white border border-gray-200 shadow-md rounded-2xl p-5 relative flex flex-col items-center">
              <div className="h-14 w-14 rounded-full bg-gradient-to-tr from-indigo-600 to-purple-600 mb-3 shadow-inner flex items-center justify-center text-xl font-bold text-white">
                U
              </div>
              <div className="h-4 w-28 bg-gray-150 rounded-full mb-1"></div>
              <div className="h-3 w-16 bg-gray-100 rounded-full mb-6"></div>
              
              <div className="w-full space-y-2">
                <div className="w-full py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-center text-xs font-semibold text-gray-700">
                  Twitter / X
                </div>
                <div className="w-full py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-center text-xs font-semibold text-gray-700">
                  GitHub Profile
                </div>
                <div className="w-full py-2.5 bg-gray-50 border border-gray-200 rounded-lg text-center text-xs font-semibold text-gray-700">
                  Personal Portfolio
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="relative max-w-7xl mx-auto px-6 py-20 border-t border-gray-200/80 z-10">
        <h2 className="text-3xl md:text-5xl font-extrabold text-center mb-4 bg-clip-text text-transparent bg-gradient-to-r from-slate-900 to-indigo-950">
          Everything You Need in One Hub
        </h2>
        <p className="text-gray-500 text-center max-w-xl mx-auto mb-16">
          Say goodbye to cluttered bios. LinkPulse powers your landing page with advanced developer-first tooling.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* Card 1 */}
          <div className="bg-white border border-gray-200/80 rounded-2xl p-6 hover:border-indigo-500/40 hover:shadow-md transition-all duration-300 shadow-sm">
            <div className="h-12 w-12 rounded-xl bg-indigo-50 border border-indigo-100 text-indigo-600 flex items-center justify-center mb-6">
              <Link2 className="h-6 w-6" />
            </div>
            <h3 className="text-xl font-bold mb-3 text-gray-900">Consolidate Everything</h3>
            <p className="text-gray-500 text-sm leading-relaxed">
              Add links to your social profiles, developer portfolios, articles, repositories, or digital storefronts. Present them beautifully in one place.
            </p>
          </div>

          {/* Card 2 */}
          <div className="bg-white border border-gray-200/80 rounded-2xl p-6 hover:border-indigo-500/40 hover:shadow-md transition-all duration-300 shadow-sm">
            <div className="h-12 w-12 rounded-xl bg-indigo-50 border border-indigo-100 text-indigo-600 flex items-center justify-center mb-6">
              <BarChart3 className="h-6 w-6" />
            </div>
            <h3 className="text-xl font-bold mb-3 text-gray-900">Real-Time Analytics</h3>
            <p className="text-gray-500 text-sm leading-relaxed">
              Track total clicks, custom referrers, traffic location logs, and daily charts. Understand exactly where your audience is coming from.
            </p>
          </div>

          {/* Card 3 */}
          <div className="bg-white border border-gray-200/80 rounded-2xl p-6 hover:border-indigo-500/40 hover:shadow-md transition-all duration-300 shadow-sm">
            <div className="h-12 w-12 rounded-xl bg-indigo-50 border border-indigo-100 text-indigo-600 flex items-center justify-center mb-6">
              <Zap className="h-6 w-6" />
            </div>
            <h3 className="text-xl font-bold mb-3 text-gray-900">Blazing Fast Go Backend</h3>
            <p className="text-gray-500 text-sm leading-relaxed">
              Leverages an ultra-optimized Go backend for redirection. Links load instantly, maximizing conversion rates and keeping users engaged.
            </p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="relative max-w-7xl mx-auto px-6 py-10 border-t border-gray-200/50 text-center text-sm text-gray-400 z-10">
        &copy; {new Date().getFullYear()} LinkPulse. All rights reserved. Built with Go and React.
      </footer>
    </div>
  );
}
